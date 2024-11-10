package wallet

type WalletRepository interface {
	NewExpense(input TransactionModel) error
	NewIncome(input TransactionModel) error
}
