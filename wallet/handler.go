package wallet

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type walletHandler struct {
	logger *zap.Logger
	// walletService WalletService
	db *sql.DB
}

func NewWalletHandler(logger *zap.Logger, db *sql.DB) *walletHandler {
	return &walletHandler{
		logger: logger,
		// walletService: walletService,
		db: db,
	}
}

func (r *walletHandler) NewTransaction(c *gin.Context) {
	var Body TransactionModel
	if err := c.ShouldBindJSON(&Body); err != nil {
		r.logger.Error("invalid request body", zap.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	usr := c.MustGet("user").(*UserModel)

	// * 1000 for not floating

	sql := `INSERT INTO transactions (user_id, category_id, amount, transaction_type) VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(sql, usr.UserID, Body.Category, Body.Amount*1000, Body.Type)
	if err != nil {
		r.logger.Error("failed to insert new expense", zap.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "failed to insert new expense"})
		return
	}

	c.JSON(200, gin.H{"message": "success"})
	return
}

func (r *walletHandler) GetSummary(c *gin.Context) {
	usr := c.MustGet("user").(*UserModel)

	sql := `select sum(amount), transaction_type from transactions where user_id = $1 group by transaction_type`
	rows, err := r.db.Query(sql, usr.UserID)
	if err != nil {
		r.logger.Error("failed to get summary", zap.String("error", err.Error()))
		c.JSON(500, gin.H{"error": "failed to get summary"})
		return
	}

	var income, expense int64 = 0, 0
	for rows.Next() {
		var amount int64
		var transactionType string
		err := rows.Scan(&amount, &transactionType)
		if err != nil {
			r.logger.Error("failed to scan row", zap.String("error", err.Error()))
			c.JSON(500, gin.H{"error": "failed to scan row"})
			return
		}
		if transactionType == "Income" {
			income += amount
		} else {
			expense -= amount
		}
	}

	c.JSON(200, gin.H{"income": income / 1000, "expense": expense / 1000})
	return
}

// func (r *walletHandler) NewIncome(c *gin.Context) {
// 	var Body IncomeModel
// 	if err := c.ShouldBindJSON(&Body); err != nil {
// 		r.logger.Error("invalid request body", zap.String("error", err.Error()))
// 		c.JSON(400, gin.H{"error": "invalid request body"})
// 		return
// 	}

// 	usr := c.MustGet("user").(UserModel)

// 	sql := `INSERT INTO transactions (user_id, category_id, amount, transaction_type) VALUES ($1,$2,$3,$4)`
// 	_, err := r.db.Exec(sql, usr.UserID, Body.Category, Body.Amount, "INCOME")
// 	if err != nil {
// 		r.logger.Error("failed to insert new expense", zap.String("error", err.Error()))
// 		c.JSON(500, gin.H{"error": "failed to insert new expense"})
// 		return
// 	}

// 	c.JSON(200, gin.H{"message": "success"})
// 	return
// }
