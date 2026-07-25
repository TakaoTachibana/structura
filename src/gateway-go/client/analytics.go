package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type AnalyticsClient struct {
	rBaseURL string
	juliaBaseURL string
	httpClient *http.Client
}

func NewAnalyticsClient(rBaseURL, juliaBaseURL string) *AnalyticsClient {
	return &AnalyticsClient {
		rBaseURL: rBaseURL,
		juliaBaseURL: juliaBaseURL,
		httpClient: &http.Client {
			Timeout: 5 * time.Second,
		},
	}
}

type MetricPoint struct {
	Time float64 `json:"time"`
	Value float64 `json:"value"`
}

type RAnalyzeRawResponse struct {
	Status []string `json:"status"`
	IsPhaseTransition []bool `json:"is_phase_transition"`
	Metrics struct {
		Variance []float64 `json:"variance"`
		Autocorrelation []float64 `json:"autocorrelation"`
		Criticality []float64 `json:"criticality_socre"`
	} `json:"metrics"`
	Timestamp []string `json:"timestamp"`
}

type RAnalyzeResult struct {
	Status string `json:"status"`
	IsPhaseTransition bool `json:"is_phase_transition"`
	Variance float64 `json:"variance"`
	Autocorrelation float64 `json:"autocorrelation"`
	CriticalityScore float64 `json:"criticality_score"`
	Timestamp string `json:"timestamp"`
}

type JuliaSolveRequest struct {
	InitialQueue float64 `json:"initial_queue"`
	ArraivalRate float64 `json:"arrival_rate"`
	ServiceRate float64 `json:"service_rate"`
	Horizon float64 `json:"horizon"`
}

type JuliaSolveResponse struct {
	Status string `json:"status"`
	TimeSteps []float64 `json:"time_steps"`
	PredictedQueue []float64 `json:"predicted_queue"`
	terminalQueue float64 `json:"terminal_queue"`
	OverflowRisk bool `json:"overflow_risk"`
}

type CombinedEvaluationResponse struct {
	PhaseTransition RAnalyzeResult `json:"phase_transition"`
	QueuePrediction JuliaSolveResponse `json:"queue_prediction"`
	EvaluatedAt time.Time `json:"evaluated_at"`
}

func (c *AnalyticsClient) EvaluateAll(
	ctx context.Context,
	metrics []MetricPoint,
	juliaReq JuliaSolveRequest,
) (*CombinedEvaluationResponse, error) {
	g, ctx :=errgroup.WithContext(ctx)

	var rRaw RAnalyzeRawResponse
	var juliaRes JuliaSolveResponse

	g.Go(func() error {
		bodyBytes, err := json.Marshal(metrics)
		if err != nil {
			return fmt.Errorf("R request marshal error: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rBaseURL+"/analyze", bytes.NewBuffer(bodyBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err  != nil {
			return fmt.Errorf("R service error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("R service retrun status %d", resp.StatusCode)
		}

		return json.NewDecoder(resp.Body).Decode(&rRaw)
	})

	g.Go(func() error {
		bodyBytes, err := json.Marshal(juliaReq)
		if err != nil {
			return fmt.Errorf("Julia request marshal error: w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.juliaBaseURL+"/solve", bytes.NewBuffer(bodyBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Context-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("Julia service error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Julia service  returned  status %d", resp.StatusCode)
		}

		return json.NewDecoder(resp.Body).Decode(&juliaRes)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	rResult := RAnalyzeResult {
		Status: getFirst(rRaw.Status, "UNKNOWN"),
		IsPhaseTransition: getFirst(rRaw.IsPhaseTransition, false),
		Variance: getFirst(rRaw.Metrics.Variance, 0.0),
		Autocorrelation: getFirst(rRaw.Metrics.Autocorrelation, 0.0),
		CriticalityScore: getFirst(rRaw.Metrics.Criticality, 0.0),
		Timestamp: getFirst(rRaw.Timestamp, ""),
	}

	return &CombinedEvaluationResponse {
		PhaseTransition: rResult,
		QueuePrediction: juliaRes,
		EvaluatedAt: time.Now().UTC(),
	}, nil
}

func getFirst[T any](slice []T, fallback T) T {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}



	
