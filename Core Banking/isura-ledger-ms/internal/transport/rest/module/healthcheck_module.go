package module

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/handler"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/types"
)

type HealthCheck struct {
	pg          *postgres.Postgres
	serviceName string
}

func NewHealthCheck(pg *postgres.Postgres, serviceName string) *HealthCheck {
	return &HealthCheck{
		pg:          pg,
		serviceName: serviceName,
	}
}

func (r *HealthCheck) Routes() types.RouteType {
	return types.RouteType{
		{
			Method:      http.MethodGet,
			Path:        "/health",
			Handler:     handler.HealthCheck(r.pg, r.serviceName),
			Description: "Health Check",
			Middlewares: types.Middlewares{},
		},
	}
}
