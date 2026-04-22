package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/infra/composition"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
)

func main() {
	deps := server.BuildBaseDeps()

	composer := composition.NewComposer(deps)
	grpcSrv := composer.GRPCServer()

	httpSrv := server.NewHTTPServer(*deps)

	go httpSrv.Start()
	go grpcSrv.Start()

	gracefulShutdown(grpcSrv, httpSrv, *deps)
}

func gracefulShutdown(grpcSrv *server.GRPCServer, httpSrv *server.HTTPServer, deps server.BaseDeps) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	deps.Log.InfoText("isura-ledger-ms shutting down...")

	// 1. para de aceitar novos requests gRPC
	grpcSrv.Stop()
	deps.Log.InfoText("Stop grpc server")

	// 2. para o HTTP com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Stop(ctx)
	deps.Log.InfoText("Stop http server")

	// 3. fecha infraestrutura — só depois que os servers pararam
	deps.Pg.Close()
	deps.Log.InfoText("Close connection postgres")

	deps.Prom.Close()
	deps.Log.InfoText("Close connection prometheus")

	deps.Log.InfoText("shutdown complete")
	os.Exit(0)
}
