package server

import (
	"fmt"
	"log/slog"
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
		badRequest(c, "")
		return
	}

	order := &models.Order{
		NumberOfItems: orderRequest.NumberOfItems,
	}
	err = s.orderService.CreateOrder(c.Request.Context(), order)
	if err != nil {
		internalServerError(c)
		return
	}

	created(c, "order created successfully", order)
}

func (s *Server) GetOrderHandler(c *gin.Context) {
	orderIDString := c.Param("id")
	orderID, err := strconv.Atoi(orderIDString)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	order, err := s.db.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error getting order: %v", err))
		notFound(c)
		return
	}

	ok(c, "successful", order)
}

func (s *Server) GetAllOrdersHandler(c *gin.Context) {
	shipping, err := s.db.GetOrdersShipping(c.Request.Context())
	if err != nil {
		errMessage := fmt.Sprintf("error getting order: %v", err)
		s.logger.Error(errMessage)
		internalServerError(c)
	}

	ok(c, "successful", shipping)
}
