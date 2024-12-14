package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupCorsConfig(r *gin.Engine) {
	if s.config.FrontendURL != "" {
		corsConfig := cors.Config{
			AllowOrigins:     []string{s.config.FrontendURL},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
		r.Use(cors.New(corsConfig))
	}

}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	s.setupCorsConfig(r)

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/orders", s.CreateOrderHandler)
	r.GET("/orders/:id", s.GetOrderHandler)

	r.GET("/orders", s.GetAllOrdersHandler)

	return r
}
