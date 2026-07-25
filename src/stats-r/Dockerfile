FROM r-base:latest

RUN apt-get update && apt-get install -y \
		r-cran-plumber \
		r-cran-mgcv \
		r-cran-jsonlite \
		&& rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

EXPOSE 8000

CMD ["Rscript", "-e", "library(plumber); pr('src/main.R') |> pr_run(host='0.0.0.0', port=8000)"]

