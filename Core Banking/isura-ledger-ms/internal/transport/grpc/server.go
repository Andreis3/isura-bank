package grpc

import (
	"context"

	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
)

type LedgerServer struct {
	pb.UnimplementedLedgerServiceServer
	createAccount *CreateAccountHandler
}

func NewLedgerServer(createAccount *CreateAccountHandler) *LedgerServer {
	return &LedgerServer{
		createAccount: createAccount,
	}
}

func (s *LedgerServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return s.createAccount.Handle(ctx, req)
}
