package wallet

import (
	"database/sql"
	"time"

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

	now := time.Now()
	var startDate time.Time
	var endDate time.Time
	if now.Day() >= 23 {
		startDate = time.Date(now.Year(), now.Month(), 23, 0, 0, 0, 0, time.UTC)
		nextMonth := now.AddDate(0, 1, 0)
		endDate = time.Date(nextMonth.Year(), nextMonth.Month(), 22, 23, 59, 59, 999999999, time.UTC)
	} else {
		lastMonth := now.AddDate(0, -1, 0)
		startDate = time.Date(lastMonth.Year(), lastMonth.Month(), 23, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(now.Year(), now.Month(), 22, 23, 59, 59, 999999999, time.UTC)
	}

	sql := `select sum(amount), transaction_type from transactions where user_id = $1 and transaction_date BETWEEN $2 AND $3 group by transaction_type`
	rows, err := r.db.Query(sql, usr.UserID, startDate, endDate)
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
			expense += amount
		}
	}

	c.JSON(200, gin.H{"income": income / 1000, "expense": expense / 1000, "balance": income/1000 - expense/1000})
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

// func (r *walletHandler) Reset()  {

// }
