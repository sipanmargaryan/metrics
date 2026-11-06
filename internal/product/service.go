package product

import (
	"context"
	"fmt"
	"strings"

	"metrics/internal/config"
)

type ProductService interface {
	CreateProduct(ctx context.Context, request CreateProductRequest) (*Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
}

type productService struct {
	repo   ProductRepository
	config *config.Config
}

func InitProductService(repo ProductRepository, config *config.Config) ProductService {
	return &productService{repo: repo, config: config}
}

func (service *productService) CreateProduct(ctx context.Context, request CreateProductRequest) (*Product, error) {
	name := strings.TrimSpace(request.Name)
	product, err := service.repo.Create(ctx, &Product{
		Name:              name,
		Price:             request.Price,
		InitialQuantity:   request.Quantity,
		RemainingQuantity: request.Quantity,
		WarehouseID:       request.WarehouseID,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return product, nil
}

func (service *productService) ListProducts(ctx context.Context) ([]Product, error) {
	products, err := service.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return products, nil
}
