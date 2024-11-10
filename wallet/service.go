package wallet

import (
	"go.uber.org/zap"
)

type walletSvc struct {
	logger *zap.Logger
	svc    WalletRepository
}

func NewWalletSvc(logger *zap.Logger, svc WalletRepository) WalletService {
	return &walletSvc{
		logger: logger,
		svc:    svc,
	}
}

func (r *walletSvc) NewExpense(input TransactionModel) error {
	if err := r.svc.NewExpense(input); err != nil {
		return err
	}
	return nil
}

func (r *walletSvc) NewIncome(input TransactionModel) error {
	if err := r.svc.NewIncome(input); err != nil {
		return err
	}
	return nil
}
