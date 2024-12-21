package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	adapter "github.com/eyo-chen/expense-tracker-go/internal/adapter"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase/monthlytrans"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context) error {
	logger.Register()

	logger.Info("Connecting to database...")
	mysqlDB, err := newMysqlDB()
	if err != nil {
		logger.Error("Unable to connect to mysql database", "error", err)
		return err
	}
	defer mysqlDB.Close()

	// Setup adapter and usecase
	adapter := adapter.New(mysqlDB, nil, nil, nil, nil, "", "")
	monthlyTransUC := monthlytrans.New(adapter.MonthlyTrans, adapter.Transaction)

	// Execute monthly transaction aggregation for previous month
	previousMonth := time.Now().AddDate(0, -1, 0)
	if err := monthlyTransUC.Create(ctx, previousMonth); err != nil {
		logger.Error("Failed to create monthly transaction aggregation", "error", err)
		return err
	}

	logger.Info("Successfully created monthly transaction aggregation", "month", previousMonth.Format(time.DateOnly))
	return nil
}

func newMysqlDB() (*sql.DB, error) {
	config := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"name":     os.Getenv("DB_NAME"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config["user"], config["password"], config["host"], config["port"], config["name"])
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
