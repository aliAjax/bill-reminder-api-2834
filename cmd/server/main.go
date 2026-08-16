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

	"bill-reminder-api/internal/bill"
	"bill-reminder-api/internal/config"
	"bill-reminder-api/internal/httpapi"
	"bill-reminder-api/internal/reminder"
	"bill-reminder-api/internal/statistics"
	storepkg "bill-reminder-api/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	jsonStore, err := storepkg.NewJSONStore(cfg.DataFile)
	if err != nil {
		log.Fatalf("init json store: %v", err)
	}

	billRepository := bill.NewRepository(jsonStore)
	reminderRepository := reminder.NewRepository(jsonStore)
	statisticsRepository := statistics.NewRepository(jsonStore)

	billService := bill.NewService(billRepository)
	reminderService := reminder.NewService(reminderRepository)
	statisticsService := statistics.NewService(statisticsRepository)

	billHandler := bill.NewHandler(billService)
	reminderHandler := reminder.NewHandler(reminderService)
	statisticsHandler := statistics.NewHandler(statisticsService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpapi.Success(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	billHandler.RegisterRoutes(mux)
	reminderHandler.RegisterRoutes(mux)
	statisticsHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("bill-reminder-api listening on :%s, data file: %s", cfg.Port, cfg.DataFile)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-shutdownSignal.Done()
	log.Println("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
