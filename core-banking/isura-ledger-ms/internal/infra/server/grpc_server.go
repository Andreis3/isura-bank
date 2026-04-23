package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcTransport "github.com/andreis3/isura-ledger-ms/internal/transport/grpc"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/interceptor"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	deps       *BaseDeps
}

func NewGRPCServer(
	deps *BaseDeps,
	ledgerServer *grpcTransport.LedgerServer,
) *GRPCServer {
	start := time.Now()

	// gRPC server com interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor(deps.Log.SlogJSON()),
			interceptor.MetricsInterceptor(deps.Prom),
			interceptor.TracingInterceptor(deps.Tracer),
		),
	)

	// registra todos os módulos
	registry := grpcTransport.NewServerRegistry(grpcServer, grpcTransport.NewLedgerModule(ledgerServer))
	registry.RegisterAll()
	reflection.Register(grpcServer)

	deps.Log.InfoText("ledger-svc started",
		slog.String("grpc_port", deps.Cfg.Servers.GRPC.Port),
		slog.String("metrics_port", deps.Cfg.Servers.HTTP.Port),
		slog.String("startup_time", time.Since(start).String()),
	)

	return &GRPCServer{
		deps:       deps,
		grpcServer: grpcServer,
	}
}

func (s *GRPCServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.deps.Cfg.Servers.GRPC.Port))
	if err != nil {
		s.deps.Log.CriticalText("grpc server failed to listen",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		s.deps.Log.CriticalText("grpc server failed to serve",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func (s *GRPCServer) Stop() {
	s.grpcServer.GracefulStop()
}
