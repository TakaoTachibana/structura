package handler

import (
	"encoding/json"
	"net/http"

	"gateway-go/client"
)

type AnalyticsHandler struct {
	client *client.AnalyticsClient
}

func NewAnalyticsHandler(c *client.AnalyticsClient) *AnalyticsHandler {
	return &AnalyticsHandler{client: c}
}

type EvalRequestPayload struct {
	Metrics []client.MetricPoint `json:"metrics"`
	SolverParams client.JuliaSolveRequest `json:"solver_params"`
}

func (h *AnalyticsHandler) HandleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload EvalRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body: "+ err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.client.EvaluateAll(r.Context(), payload.Metrics, payload.SolverParams)
	if err != nil {
		http.Error(w, "Analytics comutation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "appliaction/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

