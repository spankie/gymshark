package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spankie/gymshark/config"
	"github.com/spankie/gymshark/database"
	"github.com/spankie/gymshark/services"
)

func setupHTTPServer(t *testing.T, conf *config.Configuration, dbService database.Service) {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	orderService := services.NewOrderService(dbService, logger)

	server := NewServer(conf, dbService, orderService, logger)
	httpServer := server.NewHTTPServer()
	if httpServer == nil {
		t.Error("server creation failed")
	}

	// Start the server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("server error: %v", err)
		}
	}()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shutdown server: %v", err)
		}
	})
}

func createDBAndHTTPServer(t *testing.T, conf *config.Configuration) database.Service {
	t.Helper()
	database.CreatePostgresDBContainer(t, conf)
	dbService, err := database.NewPostgresDBService(
		conf.DbPort,
		conf.EnableDBSSL,
		conf.DbHost,
		conf.DbUsername,
		conf.DbPassword,
		conf.DbName,
	)
	if err != nil {
		t.Fatalf("error creating database service: %v", err)
	}
	setupHTTPServer(t, conf, dbService)
	return dbService
}

// Test that server is created successfully and can receive requests
func TestCreateNewHttpServer(t *testing.T) {
	conf := &config.Configuration{
		Port:       "8080",
		DbHost:     "localhost",
		DbPort:     5432,
		DbUsername: "spankie",
		DbPassword: "spankie",
		DbName:     "gymshark",
	}

	createDBAndHTTPServer(t, conf)

	// Make a request to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s", conf.Port))
	if err != nil {
		t.Errorf("failed to make request to server: %v", err)
	}

	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}
