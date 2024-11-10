package wallet

// type ExpenseModel struct {
// 	Amount   float64 `json:"amount"`
// 	Category string  `json:"category"`
// 	Type     string  `json:"type"`
// }

// type IncomeModel struct {
// 	Amount   float64 `json:"amount"`
// 	Category string  `json:"category"`
// 	Type     string  `json:"type"`
// }

type TransactionModel struct {
	Amount   float64 `json:"amount,required"`
	Category int8    `json:"category,required"`
	Type     string  `json:"type,required"`
}

type UserModel struct {
	UserID   int    `json:"user_id" db:"user_id"`   // Unique identifier for the user
	Username string `json:"username" db:"username"` // Username for the user
	Password string `json:"password" db:"password"` // Hashed password for the user
	Name     string `json:"name" db:"name"`         // Email address for the user
}
