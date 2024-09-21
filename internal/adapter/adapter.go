package model

import (
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/icon"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/maincateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/subcateg"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/transaction"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/user"
	"github.com/eyo-chen/expense-tracker-go/internal/adapter/repository/usericon"
	redisservice "github.com/eyo-chen/expense-tracker-go/internal/adapter/service/redis"
	s3service "github.com/eyo-chen/expense-tracker-go/internal/adapter/service/s3"
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
	S3Service    *s3service.Service
}

func New(mysqlDB *sql.DB,
	redisClient *redis.Client,
	s3Client *s3.Client,
	presignClient *s3.PresignClient,
	bucket string,
) *Adapter {
	return &Adapter{
		User:         user.New(mysqlDB),
		MainCateg:    maincateg.New(mysqlDB),
		SubCateg:     subcateg.New(mysqlDB),
		Icon:         icon.New(mysqlDB),
		Transaction:  transaction.New(mysqlDB),
		RedisService: redisservice.New(redisClient),
		UserIcon:     usericon.New(mysqlDB),
		S3Service:    s3service.New(bucket, s3Client, presignClient),
	}
}
