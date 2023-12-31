package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/OYE0303/expense-tracker-go/internal/handler"
	"github.com/OYE0303/expense-tracker-go/internal/model"
	"github.com/OYE0303/expense-tracker-go/internal/router"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger.Register()
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file", "error", err)
	}

	logger.Info("Connecting to database...")
	mysqlDB, err := newMysqlDB()
	if err != nil {
		logger.Fatal("Unable to connect to mysql database", "error", err)
	}
	defer mysqlDB.Close()
	mongoClient, err := newMongoDB()
	if err != nil {
		logger.Fatal("Unable to connect to mongo database", "error", err)
	}
	defer mongoClient.Disconnect(context.Background())
	mongoDB := mongoClient.Database(os.Getenv("MONGODB_DATABASE"))

	// Setup model, usecase, and handler
	model := model.New(mysqlDB, mongoDB)
	usecase := usecase.New(&model.User, &model.MainCateg, &model.SubCateg, &model.Icon, &model.Transaction)
	handler := handler.New(&usecase.User, &usecase.MainCateg, &usecase.SubCateg, &usecase.Transaction)
	if err := initServe(handler); err != nil {
		logger.Fatal("Unable to start server", "error", err)
	}
}

func newMongoDB() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI is not set")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newMysqlDB() (*sql.DB, error) {
	config := map[string]string{
		"name":     os.Getenv("DB_NAME"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:54321)/%s?parseTime=true", config["user"], config["password"], config["name"])
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
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
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
