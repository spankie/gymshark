package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	ok(c, "shipping Orders Api", nil)
}

func (s *Server) healthHandler(c *gin.Context) {
	_, err := s.db.Health(c.Request.Context())
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Error checking database health: %v", err))
		respondJSON(c, http.StatusServiceUnavailable, "", "service unavailable", response{
			Error: "database is down",
		})
		return
	}

	ok(c, "all systems are healthy", nil)
}
