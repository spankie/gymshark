package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	respondJSON(c, http.StatusOK, response{
		Message: "shipping Orders Api",
	})
}

func (s *Server) healthHandler(c *gin.Context) {
	_, err := s.db.Health(c.Request.Context())
	if err != nil {
		slog.Debug(fmt.Sprintf("Error checking database health: %v", err))
		respondJSON(c, http.StatusServiceUnavailable, response{
			Error: "database is down",
		})
		return
	}

	respondJSON(c, http.StatusOK, response{
		Message: "all systems are healthy",
	})
}
