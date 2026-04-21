package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/infra/configs"
	"github.com/andreis3/isura-ledger-ms/internal/infra/observability"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres"
	grpcTransport "github.com/andreis3/isura-ledger-ms/internal/transport/grpc"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/interceptor"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	prometheus *observability.Prometheus
	pool       *postgres.Postgres
	log        application.Logger
	cfg        *configs.Configs
}

func NewGRPCServer(
	cfg *configs.Configs,
	log application.Logger,
	prometheus *observability.Prometheus,
	pool *postgres.Postgres,
	ledgerServer *grpcTransport.LedgerServer,
) *GRPCServer {
	start := time.Now()

	// gRPC server com interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor(log.SlogJSON()),
			interceptor.MetricsInterceptor(prometheus)),
	)

	// registra todos os módulos
	registry := grpcTransport.NewServerRegistry(grpcServer, grpcTransport.NewLedgerModule(ledgerServer))
	registry.RegisterAll()
	reflection.Register(grpcServer)

	log.InfoText("ledger-svc started",
		slog.String("grpc_port", cfg.Servers.GRPC.Port),
		slog.String("metrics_port", cfg.Servers.HTTP.Port),
		slog.String("startup_time", time.Since(start).String()),
	)

	return &GRPCServer{
		prometheus: prometheus,
		pool:       pool,
		log:        log,
		cfg:        cfg,
		grpcServer: grpcServer,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Servers.GRPC.Port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

func (s *GRPCServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	s.log.InfoText("ledger-svc shutting down...")

	// para o gRPC gracefully — espera requests em andamento terminarem
	s.grpcServer.GracefulStop()

	//fecha postgres
	s.pool.Close()

	// fecha o prometheus
	s.prometheus.Close()

	s.log.InfoText("shutdown complete")
	os.Exit(0)
}
