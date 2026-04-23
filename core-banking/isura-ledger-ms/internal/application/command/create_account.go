package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

var ErrAccountAlreadyExists = errors.New("account already exists")

type CreateAccountInput struct {
	ExternalID     string
	AccountingType string
	Currency       string
}

type CreateAccount struct {
	accountRepository account.Repository
	log               application.Logger
}

func NewCreateAccount(
	accountRepository account.Repository,
	log application.Logger,
) *CreateAccount {
	return &CreateAccount{
		accountRepository: accountRepository,
		log:               log,
	}
}

func (c *CreateAccount) Execute(ctx context.Context, input CreateAccountInput) (string, error) {
	c.log.InfoJSON("CreateAccount received request",
		slog.Any("input", MaskInput[CreateAccountInput](input)),
	)

	existing, err := c.accountRepository.FindByExternalID(ctx, input.ExternalID)
	if err != nil && !errors.Is(err, account.ErrAccountNotFound) {
		c.log.CriticalJSON("CreateAccount failed to find account by external ID",
			append([]any{slog.String("external_id", input.ExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	if existing != nil {
		domainErr := fault.Wrap(fault.CodeConflict, "account already exists", ErrAccountAlreadyExists)
		c.log.WarnJSON("CreateAccount account already exists",
			append([]any{slog.String("external_id", input.ExternalID)},
				fault.Attrs(domainErr)...)...,
		)
		return "", domainErr
	}

	accountID := account.AccountID(uuid.New().String())

	accountEntity, err := account.NewAccount(
		accountID,
		input.ExternalID,
		account.AccountType(input.AccountingType),
		money.Currency(input.Currency),
	)
	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to create account entity",
			append([]any{slog.String("external_id", input.ExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	err = c.accountRepository.Save(ctx, accountEntity)
	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to save account",
			append([]any{slog.String("external_id", input.ExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	c.log.InfoJSON("CreateAccount account created successfully",
		slog.String("account_id", string(accountID)),
		slog.String("external_id", input.ExternalID),
	)

	return string(accountID), nil
}
