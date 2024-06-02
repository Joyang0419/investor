package accounting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"tools/errorx"
	"tools/redisx"
)

type ICommand interface {
	SetAccountingLock(ctx context.Context, accountID uint64) error
	ReleaseAccountingLock(ctx context.Context, accountID uint64) error
	Withdraw(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
	Deposit(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
	Transfer(ctx context.Context, fromAccountID, toAccountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
}

type Command struct {
	mysqlDB     *gorm.DB
	redisClient *redis.Client
}

var (
	ErrBalanceNotEnough = errors.New("balance is not enough")
)

const transactionLockKeyF = "transaction-lock-accountID:%d"

const lockMaxTime = 5 * time.Second

func transactionLockKey(accountID uint64) string {
	return fmt.Sprintf(transactionLockKeyF, accountID)
}

func (c *Command) SetAccountingLock(ctx context.Context, accountID uint64) error {
	key := transactionLockKey(accountID)
	if _, err := redisx.SetLock(ctx, c.redisClient, key, true, lockMaxTime); errorx.IsErrorExist(err) {
		return fmt.Errorf("[SetAccountingLock]SetLock err: %w, accountID: %d", err, accountID)
	}

	return nil
}

func (c *Command) ReleaseAccountingLock(ctx context.Context, accountID uint64) error {
	key := transactionLockKey(accountID)
	if err := redisx.ReleaseLock(ctx, c.redisClient, key); errorx.IsErrorExist(err) {
		return fmt.Errorf("[ReleaseAccountingLock]ReleaseLock err: %w, accountID: %d", err, accountID)
	}

	return nil
}

func (c *Command) Withdraw(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error) {
	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		enough, errAtomic := c.isBalanceEnough(ctx, tx, accountID, amount)
		if errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]isBalanceEnough err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}
		if !enough {
			return fmt.Errorf("[Withdraw]isBalanceEnough: %w, accountID: %d, amount: %f", ErrBalanceNotEnough, accountID, amount)
		}

		if transactionID, errAtomic = c.createTransaction(ctx, tx, accountID, amount, withdraw, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]createTransaction err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}

		negativeAmount := -amount
		if updatedBalance, errAtomic = c.updateBalance(ctx, tx, accountID, negativeAmount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]updateBalance err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Withdraw]Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Deposit(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error) {
	atomicOperation := func(tx *gorm.DB) error {
		if transactionID, err = c.createTransaction(ctx, tx, accountID, amount, "deposit", createdAt); err != nil {
			return fmt.Errorf("[Deposit]createTransaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		if updatedBalance, err = c.updateBalance(ctx, tx, accountID, amount, createdAt); err != nil {
			return fmt.Errorf("[Deposit]updateBalance err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		return nil
	}

	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Deposit]Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Transfer(ctx context.Context, fromAccountID, toAccountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error) {
	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		enough, errAtomic := c.isBalanceEnough(ctx, tx, fromAccountID, amount)
		if errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]isBalanceEnough err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}
		if !enough {
			return fmt.Errorf("[Withdraw]isBalanceEnough: %w, fromAccountID: %d, amount: %f", ErrBalanceNotEnough, fromAccountID, amount)
		}

		if transactionID, errAtomic = c.createTransaction(ctx, tx, fromAccountID, amount, withdraw, createdAt, toAccountID); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]createTransaction err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}

		negativeAmount := -amount
		if updatedBalance, errAtomic = c.updateBalance(ctx, tx, fromAccountID, negativeAmount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]updateBalance err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}
		if _, errAtomic = c.updateBalance(ctx, tx, toAccountID, amount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[Withdraw]updateBalance err: %w, toAccountID: %d, amount: %f", errAtomic, toAccountID, amount)
		}

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Withdraw]Transaction err: %w, fromAccountID: %d, amount: %f", err, fromAccountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) isBalanceEnough(ctx context.Context, tx *gorm.DB, accountID uint64, amount float64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Command) createTransaction(ctx context.Context, tx *gorm.DB, fromAccountID uint64, amount float64, transactionType string, createdAt time.Time, toAccountID ...uint64) (transactionID uint64, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *Command) updateBalance(ctx context.Context, tx *gorm.DB, accountID uint64, amount float64, createdAt time.Time) (updatedBalance float64, err error) {
	//TODO implement me
	panic("implement me")
}

func NewCommand(mysqlDB *gorm.DB, redisClient *redis.Client) ICommand {
	return &Command{
		mysqlDB:     mysqlDB,
		redisClient: redisClient,
	}
}
