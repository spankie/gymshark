package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spankie/gymshark/config"
	"github.com/spankie/gymshark/database"
	"github.com/spankie/gymshark/server"
	"github.com/spankie/gymshark/services"
)

func gracefulShutdown(apiServer *http.Server, logger *slog.Logger) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		// fmt.Fprintf(os.Stderr, "Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}

func configureLogger(logLevel string) (*slog.Logger, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, fmt.Errorf("error setting log level: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	// set the default in case standard library log is used
	slog.SetDefault(logger)

	return logger, nil
}

func run(apiServer *http.Server, logger *slog.Logger) {
	go func() {
		logger.Info("server listening on port", "port", apiServer.Addr)
		err := apiServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	gracefulShutdown(apiServer, logger)
}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting configuration: %v", err)
		os.Exit(1)
	}

	logger, err := configureLogger(conf.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error configuring logger: %v", err)
		os.Exit(1)
	}

	dbService, err := database.NewPostgresDBService(
		conf.DbPort, conf.DbHost, conf.DbUsername,
		conf.DbPassword, conf.DbName,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating database service: %v", err)
		os.Exit(1)
	}

	orderService := services.NewOrderService(dbService, logger)
	appServer := server.NewServer(conf, dbService, orderService, logger)

	port := fmt.Sprintf(":%s", conf.Port)
	handler := appServer.RegisterRoutes()
	run(server.NewHTTPServer(port, handler), logger)
}
