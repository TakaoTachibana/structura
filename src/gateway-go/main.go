package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gateway-go/ingest"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := "127.0.0.1:8080"
	receiver := ingest.NewPacketReceiver(addr)

	fmt.Printf("[gateway-go] High-throughput UDP Gateway runnning on %s...\n", addr)
	err := receiver.Start(ctx, func(data []byte) {
		fmt.Printf("[gateway-go] Received UDP Packet (%d bytes): %s\n", len(data), string(data))
	})

	if err != nil && err != context.Canceled {
		log.Fatalf("Fatal gateway error: %v", err)
	}
}

