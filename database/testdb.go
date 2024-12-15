package database

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/spankie/gymshark/config"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func getContainerExposedPort(container testcontainers.Container, containerPort int) (int, error) {
	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", containerPort))
	if err != nil {
		return 0, fmt.Errorf("failed to create port: %w", err)
	}
	mappedPort, err := container.MappedPort(context.Background(), port)
	if err != nil {
		return 0, fmt.Errorf("failed to get mapped port: %w", err)
	}
	return mappedPort.Int(), nil
}

func CreatePostgresDBContainer(t *testing.T, conf *config.Configuration) {
	t.Helper()
	postgresContainer, err := postgres.Run(context.Background(),
		"docker.io/postgres:14.1-alpine",
		postgres.WithDatabase(conf.DbName),
		postgres.WithUsername(conf.DbUsername),
		postgres.WithPassword(conf.DbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	})

	conf.DbPort, err = getContainerExposedPort(postgresContainer, 5432)
	if err != nil {
		t.Fatalf("getting container port should return nil error but got: %v", err)
	}
}
