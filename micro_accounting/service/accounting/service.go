package accounting

import (
	"context"
	"errors"
	"fmt"

	"protos/micro_accounting"
	"tools/errorx"
	"tools/numberx"
	"tools/reflectx"
)

type Service struct {
	micro_accounting.UnimplementedAccountingServiceServer
	Query   IQuery
	Command ICommand
}

func NewService(
	query IQuery,
	command ICommand,
) micro_accounting.AccountingServiceServer {
	return &Service{
		Query:   query,
		Command: command,
	}
}

var (
	ErrNilRequest    = errors.New("nil request")
	ErrWrongAccount  = errors.New("wrong account")
	ErrInvalidAmount = errors.New("invalid amount")
)

const (
	withdraw = "withdraw"
	deposit  = "deposit"
	transfer = "transfer"
)

func (s *Service) Withdraw(ctx context.Context, request *micro_accounting.WithdrawRequest) (resp *micro_accounting.WithdrawResponse, err error) {
	if reflectx.IsNil(request) {
		return nil, ErrNilRequest
	}
	if numberx.IsLT(request.Amount, 0) {
		return nil, fmt.Errorf("%w, amount: %f", ErrInvalidAmount, request.Amount)
	}

	existed, err := s.Query.IsAccountIDsExist(ctx, request.AccountID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Withdraw]IsAccountIDsExist err: %w", err)
	}
	if !existed {
		return nil, fmt.Errorf("[Withdraw]IsAccountIDsExist: %w, accountID: %d", ErrWrongAccount, request.AccountID)
	}

	transactionID, updatedBalance, err := s.Command.Withdraw(
		ctx,
		request.AccountID,
		request.Amount,
	)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Withdraw]Withdraw err: %w, accountID: %d, amount: %f", err, request.AccountID, request.Amount)
	}

	return &micro_accounting.WithdrawResponse{
		ID:             transactionID,
		Type:           withdraw,
		Amount:         request.Amount,
		AccountID:      request.AccountID,
		CurrentBalance: updatedBalance,
	}, nil
}

func (s *Service) Deposit(ctx context.Context, request *micro_accounting.DepositRequest) (resp *micro_accounting.DepositResponse, err error) {
	if reflectx.IsNil(request) {
		return nil, ErrNilRequest
	}
	if numberx.IsLT(request.Amount, 0) {
		return nil, fmt.Errorf("%w, amount: %f", ErrInvalidAmount, request.Amount)
	}

	existed, err := s.Query.IsAccountIDsExist(ctx, request.AccountID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Deposit]IsAccountIDsExist err: %w", err)
	}
	if !existed {
		return nil, fmt.Errorf("[Deposit]IsAccountIDsExist: %w, accountID: %d", ErrWrongAccount, request.AccountID)
	}

	transactionID, updatedBalance, err := s.Command.Deposit(ctx, request.AccountID, request.Amount)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Deposit]Deposit err: %w, accountID: %d, amount: %f", err, request.AccountID, request.Amount)
	}

	return &micro_accounting.DepositResponse{
		ID:             transactionID,
		Type:           deposit,
		Amount:         request.Amount,
		AccountID:      request.AccountID,
		CurrentBalance: updatedBalance,
	}, nil
}

func (s *Service) Transfer(ctx context.Context, request *micro_accounting.TransferRequest) (*micro_accounting.TransferResponse, error) {
	if reflectx.IsNil(request) {
		return nil, ErrNilRequest
	}
	if numberx.IsLT(request.Amount, 0) {
		return nil, fmt.Errorf("%w, amount: %f", ErrInvalidAmount, request.Amount)
	}

	existed, err := s.Query.IsAccountIDsExist(ctx, request.AccountID, request.TargetAccountID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Transfer]IsAccountIDsExist err: %w", err)
	}
	if !existed {
		return nil, fmt.Errorf("[Transfer]IsAccountIDsExist: %w, accountID: %d", ErrWrongAccount, request.AccountID)
	}

	transactionID, updatedBalance, err := s.Command.Transfer(ctx, request.AccountID, request.TargetAccountID, request.Amount)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Transfer]Deposit err: %w, accountID: %d, amount: %f", err, request.AccountID, request.Amount)
	}

	return &micro_accounting.TransferResponse{
		ID:              transactionID,
		Type:            transfer,
		Amount:          request.Amount,
		AccountID:       request.AccountID,
		TargetAccountID: request.TargetAccountID,
		CurrentBalance:  updatedBalance,
	}, nil
}
