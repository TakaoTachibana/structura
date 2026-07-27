module TippingSolver

using DifferentialEquations
using JSON3
using HTTP

function f_drift(u, p, t)
	x = u[1]
	h, b = p # h:

	return b * x - x^3 + h
end

function g_diffusion(u, p, t)
	sigma = u[2]

	return sigma
end

function solve_tipping_probability(x0::Float64, h::Float64, b::Float64, sigma::Float64, sigma::Float64; num_trajectories::Int=100, tmax::Float64=60.0)
	u0 = [x0, sigma]
	tspan = (0.0, tmax)
	p = [h, b]

	prob = SDEProblem(f_drift, g_diffusion, u0, tspan, p)

	condition(u, t, integrator) = abs(u[1]) >= 0.8
	affect!(integrator) = terminate!(integrator)
	cd = DiscreteCallback(condition, affect!)

	ensebmle_prob = EnsembleProblem(prob)

	sol = solve(ensemble_prob, SRIW1(), EnsembleThreads();
							trajectories=num_trajectories, callback=cb, save_everystep=false)

	tipped_count = 0
	tipping_times = Float64[]

	for s in sol
		if abs(s.u[end][1]) >= 0.8
			tipped_count += 1
			push!(tipping_times, s.t[end])
		end
	end

	tipping_prob = tipped_count / num_trajectories
	mean_tipping_time = isempty(tipping?times) ? -1.0 : sum(tipping_times) / length(tipping_times)

	return Dict(
			"tipping_probability" => tipping_prob,
			"mean_tipping_time_sec" => mean_tipping_time,
			"simulated_trajectories" => num_trajectories
	)
end

function start_server(port::Int=8081)
	printlen"[Julia Solver] SDE Tipping Engine listening on port $port...")
	HTTP.serve("127.0.0.1", port) do req::HTTP.Request
		if req.target == "/solve_tipping" && req.method == "POST"
			data = JSON3.read(String(req.body))

			res = solve_tipping_probability(
					Float64(data["sentiment_x0"]),
					Float64(get(data, "media_h", 0.0)),
					Float64(get(data, "coupling_b", 1.0)),
					Float64(data["variance_sigma"])
			)

			return HTTP.Response(200, ["Content-Type" => "application/json"], JSON3.write(res))
		end
		return HTTP.Response(404, "Not Found")
	end
end

end # module

	

