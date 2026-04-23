package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/infra/configs"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/module"
)

type SetupDeps struct {
	Mux      *chi.Mux
	Postgres *postgres.Postgres
	Log      application.Logger
	Conf     *configs.Configs
}

func Setup(deps *SetupDeps) {
	NewRegisterRoutes(
		deps.Mux,
		deps.Log,
		BuildRoutes(deps),
	).Register()
}

func BuildRoutes(deps *SetupDeps) []ModuleRoutes {
	return []ModuleRoutes{
		module.NewHealthCheck(deps.Postgres, deps.Conf.ApplicationName),
		module.NewMetrics(),
	}
}
