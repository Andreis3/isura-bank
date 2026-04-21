package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/composition"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
)

func main() {
	deps := server.BuildBaseDeps()

	//// gRPC server
	//server := server.NewServer(deps.Config, deps.Logger, deps.Prometheus, deps.Postgres)

	composer := composition.NewComposer(deps)
	grpcSrv := composer.GRPCServer()

	go grpcSrv.GracefulShutdown()
	grpcSrv.Start()
}
