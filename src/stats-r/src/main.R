library(plumber)
library(mgcv)
library(jsonlite)

function(req) {
	raw_body <- req$postBody
	input_data <- fromJSON(raw_body)
	df <- as.data.frame(input_data)

	if (nrow(df) < 10) {
		return(list(
			status = "INSUFFICIENT_DATA"
			criticality_score = 0.0,
			is_phase_transition = FALSE,
			message = "At least 10 data points are required for GAM fitting."
		))
	}

	gam_fit <- gam(value ~ s(time, bs = "cr"), data = df)
	residual_val <- residuals(gam_fit)

	current_variance <- var(residuals_val)
	acf_res <- acf(residuals_val, plot = FALSE, lag.max = 1)
	autocorr_lag1 <- ifelse(length(acf_res$acf) > 1, acf_res$acf[2], 0)

	criticality_score <- current_variance * (1 + max(0, autocorr_lag1))
	is_transition <- (autocorr_lag1 > 0.7) && (current_variance > 1.5)

	return(list(
		status = "SUCCESS",
		metrics = list(
			variance = current_variance,
			autocorrelation = autocorr_lag1,
			criticality_score = criticality_score
		),
		is_phase_transition = as.logical(is_transition),
		timestamp = Sys.time()
	))
}

