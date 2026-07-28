package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-go/analytics"
)

func main() {
	port := getEnv("PORT", "8080")
	rURL := getEnv("R_SERVICE_URL", "http://localhost:8082")
	juliaURL := getEnv("JULIA_SERVICE_URL", "http://localhost:8081")

	orchestrator := analytics.NewAnalyticsOrchestrator(juliaURL, rURL)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/v1/analytics/eval", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload analytics.CSDAlertPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON payloadd: %v", err), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		result, err := orchestrator.AnalyzeInParallel(ctx, payload)
		if err != nil {
			log.Printf("[ERROR] Pipline execution failed: %v", err)
			http.Error(w, fmt.Sprintf("Pipline analysis failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("[ERROR] Failed to encode response: %v", err)
		}
	})

	server := &http.Server {
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[Go Gateway] Running analytics orchestrator proxy on port :%s...", port)
		log.Printf("[Go Gateway] Target Julia Solver: %s", juliaURL)
		log.Printf("[Go Gateway] Target R Causal Engine: %s", rURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] HTTP server failed to start: %v", err)
		}
	}()

	<-stop
	log.Println("[Go Gateway] Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[Go Gateway] Server stopped clean.")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

