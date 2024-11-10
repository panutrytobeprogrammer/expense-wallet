package wallet

type WalletService interface {
	NewExpense(input TransactionModel) error
	NewIncome(input TransactionModel) error
}
