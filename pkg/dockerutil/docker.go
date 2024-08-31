package dockerutil

import (
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Container is a docker container
type Container struct {
	Port     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

// RunDocker runs a docker container with the given image type
func RunDocker(imageType Image) *Container {
	var err error
	var pool *dockertest.Pool
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("dockertest.NewPool failed: %s", err))
	}

	if err := pool.Client.Ping(); err != nil {
		panic(fmt.Sprintf("pool.Client.Ping failed: %s", err))
	}

	imageInfo, ok := imageInfos[imageType]
	if !ok {
		panic(fmt.Sprintf("imageType %d not found", imageType))
	}

	var resource *dockertest.Resource
	resource, err = pool.RunWithOptions(&imageInfo.RunOptions,
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

	port := resource.GetPort(imageInfo.Port)
	if err := pool.Retry(func() error {
		return imageInfo.CheckReadyFunc(port)
	}); err != nil {
		panic(fmt.Sprintf("pool.Retry failed: %s", err))
	}

	return &Container{
		Port:     port,
		pool:     pool,
		resource: resource,
	}
}

// PurgeDocker purges the docker container
func (c *Container) PurgeDocker() {
	if err := c.pool.Purge(c.resource); err != nil {
		panic(fmt.Sprintf("pool.Purge failed: %s", err))
	}
}
