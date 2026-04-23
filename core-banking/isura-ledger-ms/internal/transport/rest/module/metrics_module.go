package module

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/types"
)

type Metrics struct{}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Routes() types.RouteType {
	return types.RouteType{
		{
			Method:      http.MethodGet,
			Path:        "/metrics",
			Handler:     promhttp.Handler(),
			Description: "Metrics Prometheus",
			Middlewares: types.Middlewares{},
		},
	}
}
