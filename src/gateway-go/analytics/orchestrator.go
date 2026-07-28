package analytics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type MetricPoint struct {
	Time float64 `json:"time"`
	Value float64 `json:"value"`
}

type CSDAlertPayload struct {
	SentimentX0 float64 `json:"sentiment_x0"`
	VarianceSigma float64 `json:"variance_sigma"`
	Metrics []MetricPoint `json:"metrics"`
	SolverParams map[string]interface{} `json:"solver_params"`
}

type JuliaTippingResult struct {
	TippingProbability float64 `json:"tipping_probability"`
	MeanTippingTimeSec float64 `json:"mean_tipping_time_sec"`
	SimulatedTrajectories int `json:"simulated_trajectories"`
}

type RCausalResult struct {
	Status string `json:"status"`
	RootCause string `json:"root_cause"`
	CausalEdges []CausalEdge `json:"causal_edges"`
}

type CausalEdge struct {
	From string `json:"from"`
	To string `json:"to"`
	Weight float64 `json:"weight"`
}

type UnifiedAnalysisResult struct {
	Timestamp time.Time `json:"timestamp"`
	TippingPrediction JuliaTippingResult `json:"tipping_prediction"`
	CausalAnalysis RCausalResult `json:"causal_analysis"`
	ProcessingTimeMs int64 `json:"processing_time_ms"`
}

type AnalyticsOrchestrator struct {
	juliaURL string
	rURL string
	client *http.Client
}

func NewAnalyticsOrchestrator(juliaURL, rURL string) *AnalyticsOrchestrator {
	return &AnalyticsOrchestrator {
		juliaURL: juliaURL,
		rURL: rURL,
		client: &http.Client {
			Timeout: 3 * time.Second,
		},
	}
}

func (o *AnalyticsOrchestrator) AnalyzeInParallel(ctx context.Context, alert CSDAlertPayload) (*UnifiedAnalysisResult, error) {
	startTime := time.Now()
	var wg sync.WaitGroup

	var juliaRes JuliaTippingResult
	var rRes RCausalResult
	var juliaErr, rErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		payload := map[string]interface{} {
			"sentiment_x0": alert.SentimentX0,
			"variance_sigma": alert.VarianceSigma,
			"coupling_b": 1.0,
			"media_h": 0.0,
		}
		juliaErr = o.postJSON(ctx, o.juliaURL+"/solve_tipping", payload, &juliaRes)
	}()

	go func() {
		defer wg.Done()
		payload := map[string]interface{} {
			"metrics": alert.Metrics,
		}
		rErr = o.postJSON(ctx, o.rURL+"/discover_causality", payload, &rRes)
	}()

	wg.Wait()

	if juliaErr != nil {
		return nil, fmt.Errorf("julia solver error: %w", juliaErr)
	}
	if rErr != nil {
		return nil, fmt.Errorf("r causal engine error: %w", rErr)
	}

	return &UnifiedAnalysisResult {
		Timestamp: time.Now(),
		TippingPrediction: juliaRes,
		CausalAnalysis: rRes,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func (o *AnalyticsOrchestrator) postJSON(ctx context.Context, url string, input interface{}, output interface{}) error {
	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d from %s", resp.StatusCode, url)
	}

	return json.NewDecoder(resp.Body).Decode(output)
}
 
