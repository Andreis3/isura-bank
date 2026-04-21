package handler

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/translator"
)

type CreateAccountHandler struct {
	useCase *command.CreateAccount
	log     application.Logger
}

func NewCreateAccountHandler(
	useCase *command.CreateAccount,
	log application.Logger) *CreateAccountHandler {
	return &CreateAccountHandler{
		useCase: useCase,
		log:     log,
	}
}

func (h *CreateAccountHandler) Handle(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	input := command.CreateAccountInput{
		ExternalID:     req.GetExternalId(),
		AccountingType: req.GetAccountingType(),
		Currency:       req.GetCurrency(),
	}

	accountID, err := h.useCase.Execute(ctx, input)
	if err != nil {
		return nil, translator.ToGRPCError(err)
	}

	return &pb.CreateAccountResponse{
		AccountId: accountID,
	}, nil
}
