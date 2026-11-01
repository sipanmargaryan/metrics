package order

import (
	"context"
	"fmt"
	"metrics/internal/config"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		request CreateOrderRequest,
	) (*Order, error)
}

type orderService struct {
	repo   OrderRepository
	config *config.Config
}

func InitOrderService(repo OrderRepository, config *config.Config) OrderService {
	return &orderService{
		repo:   repo,
		config: config,
	}
}

func (service *orderService) CreateOrder(
	ctx context.Context,
	request CreateOrderRequest,
) (*Order, error) {
	order, err := service.repo.Create(ctx, &Order{
		Quantity:  request.Quantity,
		ProductID: request.ProductID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return order, nil
}
