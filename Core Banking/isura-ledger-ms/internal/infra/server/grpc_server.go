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
			interceptor.MetricsInterceptor(deps.Prom)),
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

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.deps.Cfg.Servers.GRPC.Port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

func (s *GRPCServer) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	s.deps.Log.InfoText("ledger-svc shutting down...")

	// para o gRPC gracefully — espera requests em andamento terminarem
	s.grpcServer.GracefulStop()

	//fecha postgres
	s.deps.Pg.Close()

	// fecha o prometheus
	s.deps.Prom.Close()

	s.deps.Log.InfoText("shutdown complete")
	os.Exit(0)
}
