using DifferentialEquations
using HTTP
using JSON3

function queue_dynamics!(dq, q, p, t)
	lambda, mu = p
	dq[1] = lambda - mu * tanh(max(0.0, q[1]))
end

function handle_solve(req::HTTP.Request)
	try
		body = JSON3.read(req.body)
		q0 = Float64(get(body, :initial_queue, 10.0))
		lambda = Float64(get(body, :arrival_rate, 120.0))
		mu = Float64(get(body, :service_rate, 100.0))
		t_max = Float64(get(body, :horizon, 5.0))

		u0 = [q0]
		p = (lambda, mu)
		tspan = (0.0, t_max)

		prob = ODEProblem(queue_dynamics!, u0, tspan, p)
		sol = solve(prob, Tsit5(), reltol=1e-6, abstol=1e-6)

		times = collect(range(0.0, t_max, length=20))
		predicted_queue = [sol(t)[1] for t in times]

		response_payload = Dict(
			"status" => "SUCCESS",
			"time_steps" => times,
			"predicted_queue" => predicted_queue,
			"terminal_queue" => predicted_queue[end],
			"overflow_risk" => predicted_queue[end] > 90.0
		)

		return HTTP.Response(200, ["Content-Type" => "application/json"], JSON3.write(response_payload))
	catch e
		return HTTP.Response(400, ["Content-Type" => "application/json"], JSON3.write(Dict("error" => string(e))))
	end
end

function router(req::HTTP.Request)
	if req.target == "/solve" && req.method == "POST"
		return handle_solve(req)
	else
		return HTTP.Response(404, "Not Found")
	end
end

println("Julia Solver API listening on http://0.0.0.0:8081 ...")
HTTP.serve(router, "0.0.0.0", 8081)

