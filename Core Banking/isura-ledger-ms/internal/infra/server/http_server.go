package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/andreis3/isura-ledger-ms/internal/transport/rest"
)

type HTTPServer struct {
	server *http.Server
	deps   BaseDeps
}

func NewHTTPServer(deps BaseDeps) *HTTPServer {
	start := time.Now()

	mux := chi.NewRouter()

	rest.Setup(&rest.SetupDeps{
		Mux:      mux,
		Postgres: deps.Pg,
		Log:      deps.Log,
		Conf:     deps.Cfg,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", deps.Cfg.Servers.HTTP.Port),
		Handler: mux,
	}

	deps.Log.InfoText("http server started",
		slog.String("port", deps.Cfg.Servers.HTTP.Port),
		slog.String("startup_time", time.Since(start).String()),
	)

	return &HTTPServer{
		server: server,
		deps:   deps,
	}
}

func (s *HTTPServer) Start() {
	if err := s.server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		s.deps.Log.CriticalText("http server failed",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func (s *HTTPServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.deps.Log.InfoText("http server shutting down...")

	if err := s.server.Shutdown(ctx); err != nil {
		s.deps.Log.ErrorText("http server shutdown error",
			slog.String("error", err.Error()))
	}

	s.deps.Log.InfoText("http server shutdown complete")
}
