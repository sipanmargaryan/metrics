package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"metrics/internal/helpers"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	List(ctx context.Context) ([]Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func InitProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (repo *productRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	err := repo.db.WithContext(ctx).Create(product).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return nil, fmt.Errorf("%w: warehouse does not exist", helpers.ErrNotFound)
		case "23505":
			return nil, fmt.Errorf("%w: product name already exists in this warehouse", helpers.ErrConflict)
		}
	}
	return product, err
}

func (repo *productRepository) List(ctx context.Context) ([]Product, error) {
	products := make([]Product, 0)
	err := repo.db.WithContext(ctx).Order("id").Find(&products).Error
	return products, err
}
