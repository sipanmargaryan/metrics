package order

import (
	"context"
	"fmt"

	"metrics/internal/config"
	"metrics/internal/helpers"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		request CreateOrderRequest,
	) (*Order, error)
}

type orderService struct {
	repo     OrderRepository
	config   *config.Config
	producer helpers.EventProducer
}

func InitOrderService(repo OrderRepository, config *config.Config, producer helpers.EventProducer) OrderService {
	return &orderService{
		repo:     repo,
		config:   config,
		producer: producer,
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

	key := helpers.GenerateKafkaKey(order.ID)
	if err := service.producer.Produce(ctx, "orders", key, order); err != nil {
		return nil, fmt.Errorf("publish created order %d: %w", order.ID, err)
	}
	return order, nil
}
