package wallet

import (
	"database/sql"

	"go.uber.org/zap"
)

type walletRepo struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewWalletRepo(logger *zap.Logger, db *sql.DB) WalletRepository {
	return &walletRepo{
		logger: logger,
		db:     db,
	}
}

func (r *walletRepo) NewExpense(input TransactionModel) error {
	return nil
}

func (r *walletRepo) NewIncome(input TransactionModel) error {
	return nil
}
