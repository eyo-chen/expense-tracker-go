package main

import (
	"context"
	"os"
	"time"

	adapter "github.com/eyo-chen/expense-tracker-go/internal/adapter"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
)

func main() {
	logger.Register()

	initEnv()

	// Setup adapter, usecase, and handler
	adapter := adapter.New(nil, nil, nil, nil, "", os.Getenv("STOCK_SERVICE_URL"))
	usecase := usecase.New(adapter.User, adapter.MainCateg, adapter.SubCateg, adapter.Icon, adapter.Transaction, adapter.MonthlyTrans, adapter.RedisService, adapter.UserIcon, adapter.S3Service, adapter.StockService, adapter.HistoricalPortfolioService)

	userID := 1

	if err := usecase.HistoricalPortfolio.Create(context.Background(), int32(userID), time.Now()); err != nil {
		print("error", err)
	}

	print("Create successfully")
}

func initEnv() {
	env := os.Getenv("GO_ENV")
	if env == "development-docker" || env == "production" {
		return
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file", "error", err)
	}
}
