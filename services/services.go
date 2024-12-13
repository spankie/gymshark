package services

import (
	"context"

	"github.com/spankie/gymshark/database/models"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}
