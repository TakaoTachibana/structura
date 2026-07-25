FROM julia:1.10

WORKDIR /app

RUN julia -e 'using Pkg; Pkg.add(["DifferentialEquations", "HTTP", "JSON3"])'

COPY . .

EXPOSE 8081

CMD ["julia", "src/main.jl"]

