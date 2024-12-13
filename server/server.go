package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spankie/gymshark/config"
	"github.com/spankie/gymshark/database"
	"github.com/spankie/gymshark/services"
)

type Server struct {
	config       *config.Configuration
	db           database.Service
	orderService services.OrderService
	logger       *slog.Logger
}

type response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewServer(config *config.Configuration, dbService database.Service,
	orderService services.OrderService, logger *slog.Logger) *Server {
	NewServer := &Server{
		config:       config,
		db:           dbService,
		orderService: orderService,
		logger:       logger,
	}

	return NewServer
}

// NewHTTPServer creates a new http server instance
func NewHTTPServer(port string, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         port,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// decode reads the request body and decodes it into the provided interface
func decode(c *gin.Context, v interface{}) error {
	err := c.ShouldBindJSON(v)
	if err != nil {
		return fmt.Errorf("error decoding request: %w", err)
	}

	return nil
}

func respondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func badRequest(c *gin.Context) {
	respondJSON(c, http.StatusBadRequest, response{
		Error: "bad request",
	})
}

func internalServerError(c *gin.Context) {
	respondJSON(c, http.StatusInternalServerError, response{
		Error: "internal server error",
	})
}

func notFound(c *gin.Context) {
	respondJSON(c, http.StatusNotFound, response{
		Error: "not found",
	})
}
