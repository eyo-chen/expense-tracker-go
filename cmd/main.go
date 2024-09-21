package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	adapter "github.com/eyo-chen/expense-tracker-go/internal/adapter"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/service/s3"
	"github.com/eyo-chen/expense-tracker-go/internal/handler"
	"github.com/eyo-chen/expense-tracker-go/internal/router"
	"github.com/eyo-chen/expense-tracker-go/internal/usecase"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger.Register()

	initEnv()

	logger.Info("Connecting to database...")
	mysqlDB, err := newMysqlDB()
	if err != nil {
		logger.Fatal("Unable to connect to mysql database", "error", err)
	}
	defer mysqlDB.Close()

	logger.Info("Applying schema migrations...")
	if err := applySchemaMigrations(mysqlDB); err != nil {
		logger.Fatal("Unable to apply schema migrations", "error", err)
	}

	logger.Info("Applying data migrations...")
	if err := applyDataMigrations(mysqlDB); err != nil {
		logger.Fatal("Unable to apply data migrations", "error", err)
	}

	logger.Info("Connecting to redis...")
	redisClient, err := newRedisClient()
	if err != nil {
		logger.Fatal("Unable to connect to redis", "error", err)
	}
	defer redisClient.Close()

	s3Client, presignClient, err := s3.NewS3Clients(os.Getenv("AWS_REGION"))
	if err != nil {
		logger.Fatal("Unable to create S3 clients", "error", err)
	}

	// Setup adapter, usecase, and handler
	adapter := adapter.New(mysqlDB, redisClient, s3Client, presignClient, os.Getenv("AWS_BUCKET"))
	usecase := usecase.New(adapter.User, adapter.MainCateg, adapter.SubCateg, adapter.Icon, adapter.Transaction, adapter.RedisService, adapter.UserIcon, adapter.S3Service)
	handler := handler.New(usecase.User, usecase.MainCateg, usecase.SubCateg, usecase.Transaction, usecase.Icon, usecase.InitData)
	if err := initServe(handler); err != nil {
		logger.Fatal("Unable to start server", "error", err)
	}
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

func initServe(handler *handler.Handler) error {
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == "" {
		port = fmt.Sprintf(":%d", 8000)
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      router.New(handler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func applySchemaMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/schema/",
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func applyDataMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "data_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/data/",
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func newRedisClient() (*redis.Client, error) {
	config := map[string]string{
		"host": os.Getenv("REDIS_HOST"),
		"port": os.Getenv("REDIS_PORT"),
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config["host"], config["port"]),
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
