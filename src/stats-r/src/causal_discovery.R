library(plumber)
library(jsonlite)
library(pcalg)

#* @post /discover_causality
#* @serializer json list(auto_unbox = TRUE)
function(req) {
	raw_body <- req$postBody
	data_list <- fromJSON(raw_body)

	df <- as.data.frame(data_list$metrics)
	var_names <- colnames(df)
	num_vars <- ncol(df)

	if (nrow(df) < 20) {
		return(list(error = "Insufficient sample size for causal discovery. Need at least 20 time steps."))
	}

	suffStat <- list(C = cor(df), n = nrow(df))
	pc_fit <- pc(suffStat = suffStat, indepTest = gaussCItest, alpha = 0.05, labels = var_names, verbose = FALSE)
	adj_matrix <- as(pc_fit@graph, "matrix")

	in_degrees <- colSums(adj_matrix != 0)
	out_degrees <- rowSums(adj_matrix != 0)

	root_candidates <- names(which(out_degrees > 0 & in_degrees == 0))

	if (length(root_candidates) == 0) {
		root_cause <- names(which.max(out_degrees))
	} else {
		root_cause <- root_candidates[1]
	}

	return(list(
		status = "success",
		root_cause = root_cause,
		causal_edges = build_edge_list(adj_matrix, var_names),
		matrix = adj_matrix
	))
}

build_edge_list <- function(adj_mat, labels) {
	edges <- list()
	n <- length(labels)
	for (i in 1:n) {
		for (j in 1:n) {
			if (adj_mat[i, j] != 0) {
				edges[[length(edges) + 1]] <- list(
					from = labels[i],
					to = labels[j],
					weight = unname(adj_mat[i, j])
				)
			}
		}
	}
	return(edges)
}

