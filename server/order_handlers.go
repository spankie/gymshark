package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spankie/gymshark/database/models"
)

type CreateOrderRequest struct {
	NumberOfItems int `json:"number_of_items" binding:"required,min=1"`
}

func (s *Server) CreateOrderHandler(c *gin.Context) {
	var orderRequest CreateOrderRequest
	err := decode(c, &orderRequest)
	if err != nil {
		slog.Debug(fmt.Sprintf("Error decoding order request: %v", err))
		badRequest(c)
		return
	}

	order := &models.Order{
		NumberOfItems: orderRequest.NumberOfItems,
	}
	err = s.orderService.CreateOrder(c.Request.Context(), order)
	if err != nil {
		respondJSON(c, http.StatusInternalServerError, nil)
		return
	}

	respondJSON(c, http.StatusCreated, gin.H{"id": order.ID})
}

func (s *Server) GetOrderHandler(c *gin.Context) {
	orderIDString := c.Param("id")
	orderID, err := strconv.Atoi(orderIDString)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	order, err := s.db.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error getting order: %v", err))
		notFound(c)
	}

	respondJSON(c, http.StatusOK, order)
}

func (s *Server) GetOrderShippingHandler(c *gin.Context) {
	shipping, err := s.db.GetOrderShipping(c.Request.Context())
	if err != nil {
		errMessage := fmt.Sprintf("error getting order: %v", err)
		s.logger.Error(errMessage)
		internalServerError(c)
	}

	respondJSON(c, http.StatusOK, shipping)
}
