package server

import (
	"context"
	"fmt"
	"github.com/emersion/go-smtp"
	skybackend "github.com/kartverket/skyline/pkg/backend"
	"github.com/kartverket/skyline/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type SkylineServer struct {
	ctx     context.Context
	smtp    *smtp.Server
	metrics *http.Server
}

func NewServer(ctx context.Context, cfg *config.SkylineConfig) *SkylineServer {
	backend := skybackend.NewBackend(cfg)

	server := smtp.NewServer(backend)

	server.Addr = fmt.Sprintf(":%d", cfg.Port)
	server.Domain = cfg.Hostname
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true
	server.ErrorLog = log.Default()
	//TODO make adapter, or something
	server.Debug = os.Stdout

	return &SkylineServer{
		smtp:    server,
		ctx:     ctx,
		metrics: metricsServer(cfg.MetricsPort),
	}
}

func metricsServer(metricsPort uint) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: mux,
	}
}

func (s *SkylineServer) Serve() {
	// TODO: Use context
	// TODO: Ctrl+C / signal handler

	go func() {
		slog.Info("Starting SkylineServer at " + s.smtp.Addr)
		if err := s.smtp.ListenAndServe(); err != nil {
			slog.Error("could not start SMTP server", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		slog.Info("Serving metrics at " + s.metrics.Addr)
		if err := s.metrics.ListenAndServe(); err != nil {
			slog.Error("could not start metrics server", "error", err)
			os.Exit(1)
		}
	}()

	select {
	case <-s.ctx.Done():
		shutdownCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))

		slog.Info("shutting down")

		go func() {
			err := s.smtp.Shutdown(shutdownCtx)
			if err != nil {
				slog.Warn("could not shut down SMTP server", "error", err)
			}
		}()

		go func() {
			err := s.metrics.Shutdown(shutdownCtx)
			if err != nil {
				slog.Warn("could not shut down metrics server", "error", err)
			}
		}()
	}
}
