package gnmi

import (
	"context"
	"log"
	"math/rand"
	"time"

	"gateway-go/internal/listener"
)

type GNMIListener struct {
	running bool
}

func NewGNMIListener() *GNMIListener {
	return &GNMIListener{}
}

func (l *GNMIListener) Name() string {
	return "gNMI-gRPC-Listener"
}

func (l *GNMIListener) Start(ctx context.Context, dataCh chan<- listener.TelemetryData) error {
	l.running = true
	log.Printf("[%s] Initializing gNMI gRPC Telemetry Stream Subscriptions...", l.Name())

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] Stopping gNMI Listener...", l.Name())
				return
			case t := <-ticker.C:
				simulatedInOctets := float64(1000000 + rand.Intn(500000))

				dataCh <- listener.TelemetryData {
					Timestamp: t.UnixNano(),
					Source: "gnmi",
					Metric: "interface_in_octets_bytes",
					Value: simulatedInOctets,
					Tags: map[string]string {
						"device": "core-router-01",
						"interface": "Ethernet1/1",
					},
				}
			}
		}
	}()

	return nil
}

func (l *GNMIListener) Stop() error {
	l.running = false
	return nil
}

