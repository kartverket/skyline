package server

import (
	"context"
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/kartverket/skyline/pkg/backend"
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

func NewServer(ctx context.Context, port uint, metricsPort uint, hostname string, debug bool) *SkylineServer {
	be := backend.NewBackend("username", "password")

	s := smtp.NewServer(be)

	s.Addr = fmt.Sprintf(":%d", port)
	s.Domain = hostname
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	s.ErrorLog = log.Default()
	if debug {
		s.Debug = os.Stdout
	}

	return &SkylineServer{
		smtp:    s,
		ctx:     ctx,
		metrics: metricsServer(metricsPort),
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
