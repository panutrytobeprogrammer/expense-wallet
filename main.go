package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panutrytobeprogrammer/expense-wallet/framework"
	"github.com/panutrytobeprogrammer/expense-wallet/libs"
	"github.com/panutrytobeprogrammer/expense-wallet/middleware"
	"github.com/panutrytobeprogrammer/expense-wallet/wallet"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("failed to init log: %v\n", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("failed to sync log: %v\n", err)
		}
	}()

	postgresDb := libs.ConnectPostgreSQL(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	h := wallet.NewWalletHandler(logger, postgresDb)

	app := framework.Ginapp()
	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "app is running healthy")
	})

	AuthMiddleware := middleware.NewMiddleware(logger, postgresDb)

	app.Use(AuthMiddleware.AuthRequire)
	app.GET("/api/v1/transactions", h.GetSummary)
	app.POST("/api/v1/transactions", h.NewTransaction)

	defer postgresDb.Close()

	srv := http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		d := time.Duration(5 * time.Second)
		fmt.Printf("shutting down int %s ...", d)

		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Info("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("HTTP server ListenAndServe: " + err.Error())
		return
	}

	<-idleConnsClosed
	fmt.Println("gracefully")
}
