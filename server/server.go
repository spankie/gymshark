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
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
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

func respondJSON(c *gin.Context, status int, message, err string, data interface{}) {
	c.JSON(status, response{
		Data:    data,
		Message: message,
		Error:   err,
	})
}

func ok(c *gin.Context, message string, data interface{}) {
	respondJSON(c, http.StatusOK, message, "", data)
}

func created(c *gin.Context, message string, data interface{}) {
	respondJSON(c, http.StatusCreated, message, "", data)
}

func badRequest(c *gin.Context, err string) {
	if err == "" {
		err = "bad request"
	}
	respondJSON(c, http.StatusBadRequest, "", err, nil)
}

func internalServerError(c *gin.Context) {
	respondJSON(c, http.StatusInternalServerError, "", "bad request", nil)
}

func notFound(c *gin.Context) {
	respondJSON(c, http.StatusNotFound, "", "bad request", nil)
}
