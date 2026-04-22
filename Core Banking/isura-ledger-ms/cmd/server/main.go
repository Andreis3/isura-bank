package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/composition"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
)

func main() {
	deps := server.BuildBaseDeps()

	composer := composition.NewComposer(deps)
	grpcSrv := composer.GRPCServer()

	httpSrv := server.NewHTTPServer(*deps)

	go grpcSrv.GracefulShutdown()
	go httpSrv.Start()
	grpcSrv.Start()
}
