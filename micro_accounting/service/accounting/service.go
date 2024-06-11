package accounting

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"protos/micro_accounting"
	"tools/errorx"
	"tools/numberx"
	"tools/reflectx"
	"tools/slicex"
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

	outputID, err := toOutputID(withdraw, transactionID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Withdraw]toOutputID err: %w", err)
	}

	return &micro_accounting.WithdrawResponse{
		ID:             outputID,
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

	outputID, err := toOutputID(deposit, transactionID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Deposit]toOutputID err: %w", err)
	}

	return &micro_accounting.DepositResponse{
		ID:             outputID,
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

	outputID, err := toOutputID(transfer, transactionID)
	if errorx.IsErrorExist(err) {
		return nil, fmt.Errorf("[Transfer]toOutputID err: %w", err)
	}

	return &micro_accounting.TransferResponse{
		ID:              outputID,
		Type:            transfer,
		Amount:          request.Amount,
		AccountID:       request.AccountID,
		TargetAccountID: request.TargetAccountID,
		CurrentBalance:  updatedBalance,
	}, nil
}

var prefixStrMap = map[string]string{
	withdraw: "W",
	deposit:  "D",
	transfer: "T",
}

func toOutputID(action string, transactionID int64, paddingSize ...int) (string, error) {
	prefix, exist := prefixStrMap[action]
	if !exist {
		return "", fmt.Errorf("[toOutputID]invalid action: %s", action)
	}
	hexStr := numberx.ToCapitalHex(transactionID)

	size := 0
	if slicex.IsNotEmpty(paddingSize) {
		size = paddingSize[0]
	}
	// 計算需要填充的零的個數
	paddingLen := size - len(prefix) - 1 - len(hexStr) // 減去前綴、分隔符和十六進位字符串的長度
	if paddingLen < 0 {
		return "", fmt.Errorf("[toOutputID]padding size too small for transaction ID")
	}
	// 自動填充零
	paddedHexStr := strings.Repeat("0", paddingLen) + hexStr
	outputID := fmt.Sprintf("%s-%s", prefix, paddedHexStr)

	return outputID, nil
}
