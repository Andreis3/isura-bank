package server

import (
	"log/slog"
	"os"

	"github.com/andreis3/isura-ledger-ms/internal/infra/configs"
	"github.com/andreis3/isura-ledger-ms/internal/infra/logger"
	"github.com/andreis3/isura-ledger-ms/internal/infra/observability"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres"
)

type BaseDeps struct {
	Config     *configs.Configs
	Logger     *logger.Logger
	Prometheus *observability.Prometheus
	Postgres   *postgres.Postgres
}

func BuildBaseDeps() *BaseDeps {
	cfg := configs.LoadConfig()
	log := logger.NewLogger()
	if cfg == nil {
		log.CriticalText("failed to load config")
		os.Exit(1)
	}

	prom, err := observability.NewPrometheus()
	if err != nil {
		log.CriticalText("failed to initialize prometheus", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pg, err := postgres.NewPostgres(cfg)
	if err != nil {
		log.CriticalText("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return &BaseDeps{
		Config:     cfg,
		Logger:     log,
		Prometheus: prom,
		Postgres:   pg,
	}
}
