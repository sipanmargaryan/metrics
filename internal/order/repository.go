package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"metrics/internal/helpers"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func InitOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (repo *orderRepository) Create(
	ctx context.Context,
	order *Order,
) (*Order, error) {

	err := repo.db.
		WithContext(ctx).
		Create(order).
		Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return nil, fmt.Errorf("%w: product does not exist", helpers.ErrNotFound)
	}

	return order, err
}
