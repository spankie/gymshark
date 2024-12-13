package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	// TODO: configure the gin logger to use slog to be consistent with the
	// rest of the application
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/orders", s.CreateOrderHandler)
	r.GET("/orders/:id", s.GetOrderHandler)

	r.GET("/shipping", s.GetOrderShippingHandler)

	return r
}
