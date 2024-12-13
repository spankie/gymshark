package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spankie/gymshark/database"
	"github.com/spankie/gymshark/database/models"
)

type service struct {
	db     database.Service
	logger *slog.Logger
}

func NewOrderService(db database.Service, logger *slog.Logger) OrderService {
	return service{
		db:     db,
		logger: logger.With("name", "order_service"),
	}
}

func (s service) CreateOrder(ctx context.Context, order *models.Order) error {
	packs, err := s.db.GetAvailableShippingPacks(ctx)
	if err != nil {
		return fmt.Errorf("could not find shipping packs: %w", err)
	}

	// get a slice of only the quantity to be used in calculating the shipping packs
	packSlice := make([]int, 0, len(packs))
	for _, v := range packs {
		packSlice = append(packSlice, v.Quantity)
	}

	if len(packSlice) < 1 {
		return fmt.Errorf("no packs to ship")
	}

	shippingPacks := findOptimalPacks(packSlice, order.NumberOfItems)
	err = s.db.CreateOrder(ctx, order, getOrderShipping(shippingPacks))
	if err != nil {
		err := fmt.Errorf("Error creating order: %w", err)
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func getOrderShipping(orderShippingPacks map[int]int) []*models.OrderShipping {
	orderShipping := make([]*models.OrderShipping, 0, len(orderShippingPacks))
	for k, v := range orderShippingPacks {
		orderShipping = append(orderShipping, &models.OrderShipping{
			PackSize:             k, // this should be the id
			ShippingPackQuantity: v,
		})
	}

	return orderShipping
}
