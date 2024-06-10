package accounting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"

	"definition/kafka/transaction"
	"definition/mysql/accounts"
	"definition/mysql/balance_change_log"
	"definition/mysql/transactions"
	"tools/errorx"
	redisTools "tools/infra/redis"
	"tools/numberx"
)

type ICommand interface {
	Withdraw(ctx context.Context, accountID int64, amount float64) (transactionID int64, updatedBalance float64, err error)
	Deposit(ctx context.Context, accountID int64, amount float64) (transactionID int64, updatedBalance float64, err error)
	Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) (transactionID int64, updatedBalance float64, err error)
}

type Command struct {
	mysqlDB     *gorm.DB
	redisClient *redis.Client
	kafkaConn   *kafka.Conn
}

var (
	ErrBalanceNotEnough = errors.New("balance is not enough")
)

const transactionLockKeyF = "transaction-lock-accountID:%d"

const lockMaxTime = 5 * time.Second

func transactionLockKey(accountID int64) string {
	return fmt.Sprintf(transactionLockKeyF, accountID)
}

func (c *Command) Withdraw(ctx context.Context, accountID int64, amount float64) (transactionID int64, updatedBalance float64, err error) {
	if _, err = redisTools.SetLock(ctx, c.redisClient, transactionLockKey(accountID), true, lockMaxTime); errorx.IsErrorExist(err) {
		return 0, 0, fmt.Errorf("[Command][Withdraw]redisTools.SetLock err: %w, accountID: %d", err, accountID)
	}
	defer func() {
		if errReleaseLock := redisTools.ReleaseLock(ctx, c.redisClient, transactionLockKey(accountID)); errorx.IsErrorExist(errReleaseLock) {
			err = fmt.Errorf("[Command][Withdraw]SetAccountingLock]redisTools.ReleaseLock err: %w, accountID: %d, oriErr: %v", errReleaseLock, accountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		currentBalance, errAtomic := accounts.GetBalance(ctx, tx, accountID)
		if errAtomic != nil {
			return fmt.Errorf("[Command][Withdraw]accounts.GetBalance err: %w, accountID: %d", err, accountID)
		}

		if numberx.IsLT(currentBalance, amount) {
			return fmt.Errorf("[Command][Withdraw]accounts.GetBalance err: %w, accountID: %d, amount: %f, currentBalance: %f", ErrBalanceNotEnough, accountID, amount, currentBalance)
		}

		if transactionID, err = transactions.Create(ctx, tx, withdraw, amount, accountID, accountID); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Withdraw]transactions.Create err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		negativeAmount := -amount
		if err = accounts.UpdateBalance(ctx, tx, accountID, negativeAmount); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Withdraw]accounts.UpdateBalance err: %w, accountID: %d", err, accountID)
		}

		afterBalance := currentBalance - amount
		if err = balance_change_log.Create(ctx, tx, accountID, transactionID, currentBalance, afterBalance); err != nil {
			return fmt.Errorf("[Command][Withdraw]balance_change_log.Create err: %w, accountID: %d", err, accountID)
		}

		updatedBalance = currentBalance - amount

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Command][Withdraw]mysqlDB.Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}

	if err = transaction.WriteMessages(
		ctx,
		c.kafkaConn,
		[]transaction.Data{
			{
				ID:              transactionID,
				Type:            withdraw,
				Amount:          amount,
				AccountID:       accountID,
				TargetAccountID: accountID,
			},
		},
	); err != nil {
		return 0, 0, fmt.Errorf("[Command][Withdraw]kafka2.NewKafkaSyncProducer err: %w", err)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Deposit(ctx context.Context, accountID int64, amount float64) (transactionID int64, updatedBalance float64, err error) {
	if _, err = redisTools.SetLock(ctx, c.redisClient, transactionLockKey(accountID), true, lockMaxTime); errorx.IsErrorExist(err) {
		return 0, 0, fmt.Errorf("[Command][Deposit]redisTools.SetLock err: %w, accountID: %d", err, accountID)
	}
	defer func() {
		if errReleaseLock := redisTools.ReleaseLock(ctx, c.redisClient, transactionLockKey(accountID)); errorx.IsErrorExist(errReleaseLock) {
			err = fmt.Errorf("[Command][Deposit]redisTools.ReleaseLock err: %w, accountID: %d, oriErr: %v", errReleaseLock, accountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		if numberx.IsZero(amount) {
			return errors.New("[Command][Deposit]amount is zero")
		}

		beforeBalance, errBeforeBalance := accounts.GetBalance(ctx, tx, accountID)
		if errorx.IsErrorExist(errBeforeBalance) {
			return fmt.Errorf("[Command][Deposit]accounts.GetBalance err: %w, accountID: %d", errBeforeBalance, accountID)
		}

		if transactionID, err = transactions.Create(ctx, tx, withdraw, amount, accountID, accountID); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Deposit]transactions.Create err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		if err = accounts.UpdateBalance(ctx, tx, accountID, amount); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Deposit]accounts.UpdateBalance err: %w, accountID: %d", err, accountID)
		}

		if updatedBalance, err = accounts.GetBalance(ctx, tx, accountID); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Deposit]accounts.GetBalance err: %w, accountID: %d", err, accountID)
		}

		if err = balance_change_log.Create(ctx, tx, accountID, transactionID, beforeBalance, updatedBalance); err != nil {
			return fmt.Errorf("[Command][Deposit]balance_change_log.Create err: %w, accountID: %d", err, accountID)
		}

		return nil
	}

	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Command][Deposit]Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}
	if err = transaction.WriteMessages(
		ctx,
		c.kafkaConn,
		[]transaction.Data{
			{
				ID:              transactionID,
				Type:            deposit,
				Amount:          amount,
				AccountID:       accountID,
				TargetAccountID: accountID,
			},
		},
	); err != nil {
		return 0, 0, fmt.Errorf("[Command][Deposit]kafka2.NewKafkaSyncProducer err: %w", err)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Transfer(
	ctx context.Context,
	fromAccountID, toAccountID int64, amount float64) (
	transactionID int64, updatedBalance float64, err error,
) {
	if _, err = redisTools.SetLock(ctx, c.redisClient, transactionLockKey(fromAccountID), true, lockMaxTime); errorx.IsErrorExist(err) {
		return 0, 0, fmt.Errorf("[Command][Transfer]redisTools.SetLock err: %w, accountID: %d", err, fromAccountID)
	}
	defer func() {
		if errReleaseLock := redisTools.ReleaseLock(ctx, c.redisClient, transactionLockKey(fromAccountID)); errorx.IsErrorExist(errReleaseLock) {
			err = fmt.Errorf("[Command][Transfer]redisTools.ReleaseLock err: %w, accountID: %d, oriErr: %v", errReleaseLock, fromAccountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		currentBalance, errAtomic := accounts.GetBalance(ctx, tx, fromAccountID)
		if errAtomic != nil {
			return fmt.Errorf("[Command][Transfer]accounts.GetBalance err: %w, accountID: %d", err, fromAccountID)
		}

		if numberx.IsLT(currentBalance, amount) {
			return fmt.Errorf("[Command][Transfer]accounts.GetBalance err: %w, accountID: %d, amount: %f, currentBalance: %f", ErrBalanceNotEnough, fromAccountID, amount, currentBalance)
		}

		if transactionID, err = transactions.Create(ctx, tx, withdraw, amount, fromAccountID, toAccountID); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Transfer]transactions.Create err: %w, accountID: %d, amount: %f", err, fromAccountID, amount)
		}

		negativeAmount := -amount
		if err = accounts.UpdateBalance(ctx, tx, fromAccountID, negativeAmount); errorx.IsErrorExist(err) {
			return fmt.Errorf("[Command][Transfer]accounts.UpdateBalance err: %w, accountID: %d", err, fromAccountID)
		}

		afterBalance := currentBalance - amount
		if err = balance_change_log.Create(ctx, tx, fromAccountID, transactionID, currentBalance, afterBalance); err != nil {
			return fmt.Errorf("[Command][Transfer]balance_change_log.Create err: %w, accountID: %d", err, fromAccountID)
		}

		updatedBalance = currentBalance - amount

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Command][Transfer]mysqlDB.Transaction err: %w, fromAccountID: %d, amount: %f", err, fromAccountID, amount)
	}

	if err = transaction.WriteMessages(
		ctx,
		c.kafkaConn,
		[]transaction.Data{
			{
				ID:              transactionID,
				Type:            withdraw,
				Amount:          amount,
				AccountID:       fromAccountID,
				TargetAccountID: toAccountID,
			},
		},
	); err != nil {
		return 0, 0, fmt.Errorf("[Command][Transfer]kafka2.NewKafkaSyncProducer err: %w", err)
	}

	return transactionID, updatedBalance, nil
}

func NewCommand(mysqlDB *gorm.DB, redisClient *redis.Client, kafkaConn *kafka.Conn) ICommand {
	return &Command{
		mysqlDB:     mysqlDB,
		redisClient: redisClient,
		kafkaConn:   kafkaConn,
	}
}
