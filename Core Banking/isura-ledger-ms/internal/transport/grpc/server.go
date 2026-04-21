package grpc

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/handler"
	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
)

type Handlers map[string]any

type LedgerServer struct {
	pb.UnimplementedLedgerServiceServer
	createAccount *handler.CreateAccountHandler
}

func NewLedgerServer(createAccount *handler.CreateAccountHandler) *LedgerServer {
	return &LedgerServer{
		createAccount: createAccount,
	}
}

func (s *LedgerServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return s.createAccount.Handle(ctx, req)
}
