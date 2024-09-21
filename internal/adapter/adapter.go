package model

import (
	"database/sql"

	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/usericon"
	redisservice "github.com/eyo-chen/expense-tracker-go/internal/adapter/service/redis"
	"github.com/redis/go-redis/v9"
)

type Adapter struct {
	User         *user.Repo
	MainCateg    *maincateg.Repo
	SubCateg     *subcateg.Repo
	Icon         *icon.Repo
	Transaction  *transaction.Repo
	RedisService *redisservice.Service
	UserIcon     *usericon.Repo
}

func New(mysqlDB *sql.DB, redisClient *redis.Client) *Adapter {
	return &Adapter{
		User:         user.New(mysqlDB),
		MainCateg:    maincateg.New(mysqlDB),
		SubCateg:     subcateg.New(mysqlDB),
		Icon:         icon.New(mysqlDB),
		Transaction:  transaction.New(mysqlDB),
		RedisService: redisservice.New(redisClient),
		UserIcon:     usericon.New(mysqlDB),
	}
}
