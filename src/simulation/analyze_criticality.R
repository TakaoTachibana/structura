library(arrow)
library(tidyverse)
library(mgcv)

df <- read_parquet("fanaticism_simulation_log.parquet")
gam_model <- gam(order_parameter_m ~ s(temperature_T, k = 15), data = df)
T_grid <- seq(min(df$temperature_T), max(df$temperature_T), length.out = 1000)
eps <- 1e-5
pred_m <- predict(gam_model, newdata = data.frame(temperature_T = T_grid))
pred_m_eps <- predict(gam_model, newdata = data.frame(temperature_T = T_grid + eps))
dm_dT <- (pred_m_eps - pred_m) / eps
tc_index <- which.max(abs(dm_dT))
Tc <- T_grid[tc_index]
cat("\n============================\n")
cat(sprintf("[Completed] Detected Tc: %.4f\n", Tc))
cat("============================\n\n")

analysis_df <- data.frame(T = T_grid, m = pred_m, slope = abs(dm_dT))

p <- ggplot() +
	geom_point(data = df, aes(x = temperature_T, y = order_parameter_m), alpha = 0.3, color = "gray30") +
	geom_line(data = analysis_df, aes(x = T, y = m), color = "purple", size = 1.2) +
	geom_vline(xintercept = Tc, color = "red", linetype = "dashed", size = 1) +
	annotate("text", x = Tc + 0.15, y = 0, label = sprintf("Tc ~ %.3f", Tc), color = "red", fontface = "bold", size = 5) +
	scale_x_reverse() +
	labs(
			 title = "Non-Linear GAM Fitting & Critical Boundary (Tc) Detection",
			 subtitle = "Red dashed line indicates the inflection point where phase transition accelerates.",
			 x = "Social Stress / Temprature (T)",
			 y = "Order Parameter (m)"
	) +
	theme_minimal()

ggsave("r_criticality_analysis.png", plot = p, width = 8, height = 5, dpi = 300)
cat("Plot saved as 'r_criticality_analysis.png'.\n")

