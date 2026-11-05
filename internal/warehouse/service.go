package warehouse

import (
	"context"
	"fmt"
	"strings"

	"metrics/internal/config"
)

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, request CreateWarehouseRequest) (*Warehouse, error)
	ListWarehouses(ctx context.Context) ([]Warehouse, error)
}

type warehouseService struct {
	repo     WarehouseRepository
	config   *config.Config
}

func InitWarehouseService(repo WarehouseRepository, config *config.Config) WarehouseService {
	return &warehouseService{repo: repo, config: config}
}

func (service *warehouseService) CreateWarehouse(ctx context.Context, request CreateWarehouseRequest) (*Warehouse, error) {
	name := strings.TrimSpace(request.Name)
	warehouse, err := service.repo.Create(ctx, &Warehouse{Name: name})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return warehouse, nil
}

func (service *warehouseService) ListWarehouses(ctx context.Context) ([]Warehouse, error) {
	warehouses, err := service.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return warehouses, nil
}
