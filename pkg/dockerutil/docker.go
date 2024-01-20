package dockerutil

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	mut      = sync.Mutex{}
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func RunDocker() string {
	mut.Lock()
	defer mut.Unlock()

	var db *sql.DB
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("dockertest.NewPool failed: %s", err))
	}

	if err := pool.Client.Ping(); err != nil {
		panic(fmt.Sprintf("pool.Client.Ping failed: %s", err))
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env:        []string{"MYSQL_ROOT_PASSWORD=root"},
	},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		panic(fmt.Sprintf("pool.RunWithOptions failed: %s", err))
	}

	port := resource.GetPort("3306/tcp")
	if err := pool.Retry(func() error {
		db, err = sql.Open("mysql", fmt.Sprintf("root:root@(localhost:%s)/mysql?parseTime=true", port))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		panic(fmt.Sprintf("pool.Retry failed: %s", err))
	}

	return resource.GetPort("3306/tcp")
}

func PurgeDocker() {
	if err := pool.Purge(resource); err != nil {
		panic(fmt.Sprintf("pool.Purge failed: %s", err))
	}
}
