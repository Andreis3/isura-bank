package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/infra/logger"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository"
	grpcTransport "github.com/andreis3/isura-ledger-ms/internal/transport/grpc"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/interceptor"
	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
)

func main() {
	logg := logger.NewLogger()

	// 1. conexão com o banco
	pool, err := pgxpool.New(context.Background(), "postgres://admin:admin@localhost:5432/isura-ledger")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// 2. infraestrutura
	//unitOfWork := uow.NewUnitOfWork(pool)
	accountRepository := repository.NewAccountRepository(pool)
	//transactionRepository := repository.NewTransaction(pool)
	//outboxRepository := repository.NewOutBoxRepository(pool)

	// 3. use cases
	createAccount := command.NewCreateAccount(accountRepository, logg)
	//command.NewCreateTransaction(unitOfWork, accountRepository, transactionRepository, outboxRepository)

	// 4. handlers gRPC
	createAccountHandler := grpcTransport.NewCreateAccountHandler(createAccount, logg)

	// 5. server
	ledgerServer := grpcTransport.NewLedgerServer(createAccountHandler)

	// 6. gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor(logg.SlogJSON())),
	)
	pb.RegisterLedgerServiceServer(grpcServer, ledgerServer)
	reflection.Register(grpcServer)

	// 7. listener TCP
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logg.CriticalText(slog.LevelError.String(), err)
		os.Exit(1)
	}

	logg.InfoText("ledger-svc listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
