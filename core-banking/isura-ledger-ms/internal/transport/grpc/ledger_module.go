package grpc

import (
	"google.golang.org/grpc"

	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
)

type LedgerModule struct {
	ledgerServer *LedgerServer
}

func NewLedgerModule(server *LedgerServer) *LedgerModule {
	return &LedgerModule{
		ledgerServer: server,
	}
}

func (m *LedgerModule) Register(server *grpc.Server) {
	pb.RegisterLedgerServiceServer(server, m.ledgerServer)
}
