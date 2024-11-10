package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panutrytobeprogrammer/expense-wallet/wallet"
	"go.uber.org/zap"
)

type storage struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewMiddleware(logger *zap.Logger, db *sql.DB) *storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}

func (r *storage) AuthRequire(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		r.logger.Error("basic auth failed")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var User wallet.UserModel
	sql := "SELECT user_id, username, password, name  FROM users WHERE username = $1"

	err := r.db.QueryRow(sql, username).Scan(&User.UserID, &User.Username, &User.Password, &User.Name)
	if err != nil {
		r.logger.Error("failed to get user", zap.String("err", err.Error()))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if User.Password != password {
		r.logger.Error("password not match")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Set("user", &User)
	c.Next()
}
