package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-go/internal/listener"
	"gateway-go/internal/listener/ebpf"
	"gateway-go/internal/listener/gnmi"
	"gateway-go/internal/listener/ipfix"
)

func main() {
	log.Printf("[Go Gateway] Initializing Multi-Protocol Telemetry Pipeline...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	dataCh := make(chan listener.TelemetryData, 100)

	ebpfListener := ebpf.NewEBPFListener()
	if err := ebpfListener.Start(ctx, dataCh); err != nil {
		log.Fatalf("Failed to start %s: %v", ebpfListener.Name(), err)
	}

	gnmiListener := gnmi.NewGNMIListener()
	if err := gnmiListener.Start(ctx, dataCh); err != nil {
		log.Fatalf("Failed to start %s: %v", gnmiListener.Name(), err)
	}

	ipfixListener := ipfix.NewIPFIXListener()
	if err := ipfixListener.Start(ctx, dataCh); err != nil {
		log.Printf("Warning: Failed to start IPFIX: %v", err)
	}

	go func() {
		for data := range dataCh {
			log.Printf("[PIPELINE <- %-5s] Metric: %-25s | Value: %10.2fns | Tags: %v",
				data.Source, data.Metric, data.Value, data.Tags)
		}
	}()

	log.Printf("[Go Gateway] Pipeline Fully Active! (eBPF + gNMI + IPFIX) Press Ctrl+C to stop.")

	<-sigCh
	log.Printf("\n[Go Gateway] Shutdonw signal received. Cleaning up...")

	cancel()
	_ = ebpfListener.Stop()
	_ = gnmiListener.Stop()
	_ = ipfixListener.Stop()

	time.Sleep(200 * time.Millisecond)
	close(dataCh)

	log.Printf("[Go Gateway] Successflly shut down.")
}



