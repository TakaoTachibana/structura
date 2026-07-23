package listener

import "context"

type TelemetryData struct {
	Timestamp int64 `json:"timestamp"`
	Source string `json:"source"`
	Metric string `json:"metric"`
	Value float64 `json:"value"`
	Tags map[string]string `json:"tags"`
}

type TelemetryListener interface {
	Name() string
	Start(ctx context.Context, dataCh chan<- TelemetryData) error
	Stop() error
}

