package accounting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"

	"tools/redisx"

	"definition/db_schema/mysql"
	"tools/errorx"
)

type ICommand interface {
	Withdraw(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
	Deposit(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
	Transfer(ctx context.Context, fromAccountID, toAccountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error)
}

type Command struct {
	mysqlDB     *gorm.DB
	redisClient *redis.Client
	kafkaClient *kafka.Client
}

var (
	ErrBalanceNotEnough = errors.New("balance is not enough")
)

const transactionLockKeyF = "transaction-lock-accountID:%d"

const lockMaxTime = 5 * time.Second

func transactionLockKey(accountID uint64) string {
	return fmt.Sprintf(transactionLockKeyF, accountID)
}

func (c *Command) setAccountingLock(ctx context.Context, accountID uint64) error {
	key := transactionLockKey(accountID)
	if _, err := redisx.SetLock(ctx, c.redisClient, key, true, lockMaxTime); errorx.IsErrorExist(err) {
		return fmt.Errorf("[accounting][Command]SetAccountingLock]SetLock err: %w, accountID: %d", err, accountID)
	}

	return nil
}

func (c *Command) releaseAccountingLock(ctx context.Context, accountID uint64) error {
	key := transactionLockKey(accountID)
	if err := redisx.ReleaseLock(ctx, c.redisClient, key); errorx.IsErrorExist(err) {
		return fmt.Errorf("[accounting][Command][ReleaseAccountingLock]ReleaseLock err: %w, accountID: %d", err, accountID)
	}

	return nil
}

func (c *Command) Withdraw(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error) {
	if err = c.setAccountingLock(ctx, accountID); err != nil {
		return 0, 0, fmt.Errorf("[accounting][Command][Withdraw]setAccountingLock err: %w, accountID: %d", err, accountID)
	}
	defer func() {
		if releaseAccountingLock := c.releaseAccountingLock(ctx, accountID); releaseAccountingLock != nil {
			err = fmt.Errorf("[accounting][Command][Withdraw]releaseAccountingLock err: %w, accountID: %d, oriErr: %v", releaseAccountingLock, accountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		enough, errAtomic := c.isBalanceEnough(ctx, tx, accountID, amount)
		if errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]isBalanceEnough err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}
		if !enough {
			return fmt.Errorf("[accounting][Command][Withdraw]isBalanceEnough: %w, accountID: %d, amount: %f", ErrBalanceNotEnough, accountID, amount)
		}

		if transactionID, errAtomic = c.createTransaction(ctx, tx, accountID, amount, withdraw, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]createTransaction err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}

		negativeAmount := -amount
		if updatedBalance, errAtomic = c.updateBalance(ctx, tx, accountID, negativeAmount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]updateBalance err: %w, accountID: %d, amount: %f", errAtomic, accountID, amount)
		}

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[accounting][Command][Withdraw]Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Deposit(ctx context.Context, accountID uint64, amount float64, createdAt time.Time) (transactionID uint64, updatedBalance float64, err error) {
	if err = c.setAccountingLock(ctx, accountID); err != nil {
		return 0, 0, fmt.Errorf("[accounting][Command][Deposit]setAccountingLock err: %w, accountID: %d", err, accountID)
	}
	defer func() {
		if releaseAccountingLock := c.releaseAccountingLock(ctx, accountID); releaseAccountingLock != nil {
			err = fmt.Errorf("[accounting][Command][Deposit]releaseAccountingLock err: %w, accountID: %d, oriErr: %v", releaseAccountingLock, accountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		if transactionID, err = c.createTransaction(ctx, tx, accountID, amount, "deposit", createdAt); err != nil {
			return fmt.Errorf("[accounting][Command][Deposit]createTransaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		if updatedBalance, err = c.updateBalance(ctx, tx, accountID, amount, createdAt); err != nil {
			return fmt.Errorf("[accounting][Command][Deposit]updateBalance err: %w, accountID: %d, amount: %f", err, accountID, amount)
		}

		return nil
	}

	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[Deposit]Transaction err: %w, accountID: %d, amount: %f", err, accountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) Transfer(
	ctx context.Context,
	fromAccountID, toAccountID uint64, amount float64, createdAt time.Time) (
	transactionID uint64, updatedBalance float64, err error,
) {
	if err = c.setAccountingLock(ctx, fromAccountID); err != nil {
		return 0, 0, fmt.Errorf("[accounting][Command][Transfer]setAccountingLock err: %w, accountID: %d", err, fromAccountID)
	}
	defer func() {
		if releaseAccountingLock := c.releaseAccountingLock(ctx, fromAccountID); releaseAccountingLock != nil {
			err = fmt.Errorf("[accounting][Command][Transfer]releaseAccountingLock err: %w, accountID: %d, oriErr: %v", releaseAccountingLock, fromAccountID, err)
		}
	}()

	atomicOperation := func(tx *gorm.DB) error {
		var errAtomic error
		enough, errAtomic := c.isBalanceEnough(ctx, tx, fromAccountID, amount)
		if errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]isBalanceEnough err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}
		if !enough {
			return fmt.Errorf("[accounting][Command][Withdraw]isBalanceEnough: %w, fromAccountID: %d, amount: %f", ErrBalanceNotEnough, fromAccountID, amount)
		}

		if transactionID, errAtomic = c.createTransaction(ctx, tx, fromAccountID, amount, withdraw, createdAt, toAccountID); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]createTransaction err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}

		negativeAmount := -amount
		if updatedBalance, errAtomic = c.updateBalance(ctx, tx, fromAccountID, negativeAmount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]updateBalance err: %w, fromAccountID: %d, amount: %f", errAtomic, fromAccountID, amount)
		}

		// TODO 這一段, 我在思考，可以用「流處理」, 交給其他服務去處理, 因為轉帳的事情，我只要確認轉出錢的人有錢就好, 然後有成功扣完款要轉出錢的人，
		// TODO 至於把錢給「被轉帳人」，這似乎也不是要馬上知道的事情, 反正確定轉帳單有成立，扣完「要轉帳人的錢」就好
		if _, errAtomic = c.updateBalance(ctx, tx, toAccountID, amount, createdAt); errorx.IsErrorExist(errAtomic) {
			return fmt.Errorf("[accounting][Command][Withdraw]updateBalance err: %w, toAccountID: %d, amount: %f", errAtomic, toAccountID, amount)
		}
		// TODO 到這邊

		return nil
	}
	if err = c.mysqlDB.Transaction(atomicOperation); err != nil {
		return 0, 0, fmt.Errorf("[accounting][Command][Withdraw]Transaction err: %w, fromAccountID: %d, amount: %f", err, fromAccountID, amount)
	}

	return transactionID, updatedBalance, nil
}

func (c *Command) isBalanceEnough(ctx context.Context, tx *gorm.DB, accountID uint64, amount float64) (bool, error) {
	currentBalance, err := new(mysql.Account).GetBalance(ctx, tx, accountID)
	if err != nil {
		return false, fmt.Errorf("[accounting][Command][isBalanceEnough]GetBalance err: %w, accountID: %d", err, accountID)
	}
	return currentBalance >= amount, nil
}

func (c *Command) createTransaction(ctx context.Context, tx *gorm.DB, fromAccountID uint64, amount float64, transactionType string, createdAt time.Time, toAccountID ...uint64) (transactionID uint64, err error) {
	return new(mysql.Transaction).Create(ctx, tx, transactionType, amount, fromAccountID, toAccountID[0], createdAt)
}

func (c *Command) updateBalance(ctx context.Context, tx *gorm.DB, accountID uint64, amount float64, updatedAt time.Time) (updatedBalance float64, err error) {
	accountModel := new(mysql.Account)
	if err = accountModel.UpdateBalance(ctx, tx, accountID, amount, updatedAt); err != nil {
		return 0, fmt.Errorf("[accounting][Command][updateBalance]UpdateBalance err: %w, accountID: %d", err, accountID)
	}

	return accountModel.GetBalance(ctx, tx, accountID)
}

func NewCommand(mysqlDB *gorm.DB, redisClient *redis.Client) ICommand {
	return &Command{
		mysqlDB:     mysqlDB,
		redisClient: redisClient,
	}
}
