package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluateAll_Success(t *testing.T) {
	rServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": ["SUCCESS"],
			"is_phase_transition": [false],
			"metrics": {"variance": [0.0015], "autocorrelation": [-0.8099], "criticality_score": [0.0]},
			"timestamp": ["2026-07-25 12:00:00"]
		}`))
	}))
	defer rServer.Close()

	juliaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "SUCCESS",
			"time_steps": [0.0, 5.0],
			"predicted_queue": [10.0, 260.0],
			"terminal_queue": 260.0,
			"overflow_risk": true
		}`))
	}))
	defer juliaServer.Close()

	client := NewAnalyticsClient(rServer.URL, juliaServer.URL)
	res, err := client.EvaluateAll(context.Background(), []MetricPoint{{Time: 1, Value: 10}}, JuliaSolveRequest{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.PhaseTransition.Status != "SUCCESS" {
		t.Errorf("Expected status SUCCESS, got %s", res.PhaseTransition.Status)
	}
}

