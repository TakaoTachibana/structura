using Pkg

# スクリプトがあるディレクトリのプロジェクト環境を自動アクティベート
Pkg.activate(@__DIR__)

# 不足しているパッケージを自動インストール
required_pkgs = ["DifferentialEquations", "StochasticDiffEq", "JSON3", "HTTP"]
for pkg in required_pkgs
	if !haskey(Pkg.project().dependencies, pkg)
		println("[Julia Setup] Installing missing package: $pkg ...")
		Pkg.add(pkg)
	end
end

module TippingSolver

using DifferentialEquations
using StochasticDiffEq
using JSON3
using HTTP

# In-place 形式のドリフト項（アロケーションゼロ）
function f_drift!(du, u, p, t)
	h, b, sigma = p
	du[1] = b * u[1] - u[1]^3 + h
	return nothing
end

# In-place 形式のディフュージョン項（加法的ノイズ）
function g_diffusion!(du, u, p, t)
	h, b, sigma = p
	du[1] = sigma
	return nothing
end

function solve_tipping_probability(x0::Float64, h::Float64, b::Float64, sigma::Float64; num_trajectories::Int=100, tmax::Float64=60.0)
	u0 = [x0] # 1次元ベクトル形式
	tspan = (0.0, tmax)
	p = (h, b, sigma)

	prob = SDEProblem(f_drift!, g_diffusion!, u0, tspan, p)

	# 閾値判定用コールバック
	condition = (u, t, integrator) -> abs(u[1]) >= 0.8
	affect! = (integrator) -> terminate!(integrator)
	cb = DiscreteCallback(condition, affect!)

	# スレッドセーフな u0 コピー関数
	prob_func = (prob, args...) -> remake(prob, u0 = copy(prob.u0))
	ensemble_prob = EnsembleProblem(prob, prob_func = prob_func)

	# SRA1 ソルバーで並列実行
	ens_sol = solve(ensemble_prob, SRA1(), EnsembleThreads();
							dt=0.05, trajectories=num_trajectories, callback=cb, save_everystep=false)

	tipped_count = 0
	tipping_times = Float64[]
	
	# ens_sol.u は各軌道 (RODESolution) の配列であることが確定しているため安全
	total_trajectories = length(ens_sol.u)

	for sim in ens_sol.u
		# sim.u[end] は最終状態の Vector{Float64} ([x_final])
		final_val = sim.u[end][1]
		final_t   = sim.t[end]

		if abs(final_val) >= 0.8
			tipped_count += 1
			push!(tipping_times, final_t)
		end
	end

	# 確率（0.0 ～ 1.0）および平均到達時間の算出
	tipping_prob = total_trajectories > 0 ? tipped_count / total_trajectories : 0.0
	mean_tipping_time = isempty(tipping_times) ? -1.0 : sum(tipping_times) / length(tipping_times)

	return Dict(
			"tipping_probability" => tipping_prob,
			"mean_tipping_time_sec" => mean_tipping_time,
			"simulated_trajectories" => total_trajectories
	)
end

function warmup()
	print("[Julia] Compiling SDE Solver (Warmup)... ")
	try
		solve_tipping_probability(0.1, 0.0, 1.0, 0.1; num_trajectories=2, tmax=5.0)
		println("Done!")
	catch e
		println("\nWarmup failed: ", e)
	end
end

function start_server(port::Int=8081)
	warmup()
	println("[Julia Solver] SDE Tipping Engine listening on port $port...")
	
	HTTP.serve("0.0.0.0", port) do req::HTTP.Request
		path = HTTP.URI(req.target).path
		clean_path = rstrip(path, '/')

		if (clean_path == "/solve" || clean_path == "/solve_tipping") && req.method == "POST"
			println("[Julia] Received request. Solving SDE...")
			try
				data = JSON3.read(String(req.body))

				x0 = Float64(get(data, "sentiment_x0", get(data, "initial_queue", 0.1)))
				h  = Float64(get(data, "media_h", 0.0))
				b  = Float64(get(data, "coupling_b", 1.0))
				s  = Float64(get(data, "variance_sigma", 0.1))

				res = solve_tipping_probability(x0, h, b, s)

				println("[Julia] Solved! Sending response.")
				return HTTP.Response(200, ["Content-Type" => "application/json"], JSON3.write(res))
			catch e
				println("[Julia Error] Exception during solve: ", e)
				return HTTP.Response(500, ["Content-Type" => "application/json"], JSON3.write(Dict("error" => string(e))))
			end
		end
		return HTTP.Response(404, ["Content-Type" => "text/plain"], "Not Found")
	end
end

end # module

TippingSolver.start_server(8081)
