package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-go/client"
	"gateway-go/handler"
)

func main() {
	port := getEnv("PORT", "8080")
	rURL := getEnv("R_SERVICE_URL", "http://localhost:8000")
	juliaURL := getEnv("JULIA_SERVICE_URL", "http://localhost:8081")

	log.Printf("[INFO] Starting Gateway Service...")
	log.Printf("[INFO] Target R Service (stats-r) : %s", rURL)
	log.Printf("[INFO] Target Julia Service (solver-julia): %s", juliaURL)

	analyticsClient := client.NewAnalyticsClient(rURL, juliaURL)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsClient)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status:":"UP"}`))
	})

	mux.HandleFunc("/api/v1/analytics/eval", analyticsHandler.HandleEvaluate)

	server := &http.Server {
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[INFO] Gateway HTTP Listener on http://0.0.0.0:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server startup failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("[INFO] Shutting down Gateway server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced shutdown: %v", err)
	}

	log.Printf("[INFO] Gateway server stopped successfully.")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

