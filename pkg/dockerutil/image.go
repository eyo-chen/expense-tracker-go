package dockerutil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
)

// Image is a type of image
type Image int32

const (
	// ImageUnspecified is a Image of type Unspecified.
	ImageUnspecified Image = iota

	// ImageMySQL is a Image of type MySQL.
	ImageMySQL

	// ImageRedis is a Image of type Redis.
	ImageRedis
)

type imageInfo struct {
	dockertest.RunOptions

	Port           string
	CheckReadyFunc func(port string) error
}

var imageInfos = map[Image]imageInfo{
	ImageMySQL: {
		RunOptions: dockertest.RunOptions{
			Repository: "mysql",
			Tag:        "8.0",
			Env:        []string{"MYSQL_ROOT_PASSWORD=root"},
		},
		Port: "3306/tcp",
		CheckReadyFunc: func(port string) error {
			db, err := sql.Open("mysql", fmt.Sprintf("root:root@(localhost:%s)/mysql?parseTime=true", port))
			if err != nil {
				return err
			}

			return db.Ping()
		},
	},
	ImageRedis: {
		RunOptions: dockertest.RunOptions{
			Repository: "redis",
			Tag:        "7.2.4",
		},
		Port: "6379/tcp",
		CheckReadyFunc: func(port string) error {
			client := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("localhost:%s", port),
			})
			ctx := context.Background()

			return client.Ping(ctx).Err()
		},
	},
}
