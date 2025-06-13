package transfer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Service struct {
	accountRepo     AccountRepository
	transactionRepo TransactionRepository
	TxManager       TxManager
	logger          *zap.Logger
}

// NewService constructs a new TransferService implementation.
func NewService(accountRepo AccountRepository, transactionRepo TransactionRepository, txManager TxManager, l *zap.Logger) *Service {
	return &Service{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		TxManager:       txManager,
		logger:          l,
	}
}

// Transfer performs a funds transfer between two accounts atomically.
func (s *Service) Transfer(ctx context.Context, sourceID, destID int, amount decimal.Decimal) (*model.Transaction, error) {
	if sourceID == destID {
		s.logger.Warn("source and destination cannot be the same")
		return nil, errors.New("cannot transfer to the same account")
	}

	var result *model.Transaction
	s.logger.Info(fmt.Sprintf("Transfer from %d to %d transaction begin", sourceID, destID))
	err := s.TxManager.Do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		sourceAcc, err := s.accountRepo.GetByID(ctx, tx, sourceID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to get source account by id %d", sourceID))
			return err
		}

		destAcc, err := s.accountRepo.GetByID(ctx, tx, destID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to get destination account by id %d", destID))
			return err
		}

		if err := sourceAcc.Withdraw(amount); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to withdraw account by id %d", sourceID))
			return err // insufficient funds, etc.
		}
		destAcc.Deposit(amount)

		// Create transaction entity
		modelTx, err := model.NewTransaction(sourceID, destID, amount)
		if err != nil {
			return err
		}

		// Persist changes
		if err := s.accountRepo.Update(ctx, tx, sourceAcc); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update source account by id %d", sourceID))
			return err
		}
		if err := s.accountRepo.Update(ctx, tx, destAcc); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update destination account by id %d", sourceID))
			return err
		}
		if err := s.transactionRepo.Save(ctx, modelTx); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to save transaction %+v", modelTx))
			return err
		}

		result = modelTx
		s.logger.Info(fmt.Sprintf("Transfer from %d to %d transaction end", sourceID, destID))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) Deposit(ctx context.Context, accountID int, amount decimal.Decimal) (*model.Account, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		s.logger.Error(fmt.Sprintf("Invalid amount %s", amount.String()))
		return nil, errors.New("amount must be positive")
	}

	var updated *model.Account

	s.logger.Info(fmt.Sprintf("Deposit from %d to %d transaction begin", accountID, amount))
	err := s.TxManager.Do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		acc, err := s.accountRepo.GetByID(ctx, tx, accountID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to get account by id %d", accountID))
			return err
		}

		acc.Deposit(amount)

		if err := s.accountRepo.Update(ctx, tx, acc); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update account by id %d", accountID))
			return err
		}

		updated = acc
		return nil
	})

	return updated, err
}

func (s *Service) Withdraw(ctx context.Context, accountID int, amount decimal.Decimal) (*model.Account, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		s.logger.Error(fmt.Sprintf("Invalid amount %s", amount.String()))
		return nil, errors.New("amount must be positive")
	}

	var updated *model.Account
	s.logger.Info(fmt.Sprintf("Withdraw from %d to %d transaction begin", accountID, amount))
	err := s.TxManager.Do(ctx, func(ctx context.Context, tx *sql.Tx) error {
		acc, err := s.accountRepo.GetByID(ctx, tx, accountID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to get account by id %d", accountID))
			return err
		}

		if err := acc.Withdraw(amount); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to withdraw account by id %d", accountID))
			return err
		}

		if err := s.accountRepo.Update(ctx, tx, acc); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update account by id %d", accountID))
			return err
		}

		updated = acc
		return nil
	})

	return updated, err
}
