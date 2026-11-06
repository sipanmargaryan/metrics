package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"metrics/internal/helpers"
	"metrics/internal/product"
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
	if order.Quantity <= 0 {
		return nil, fmt.Errorf("%w: order quantity must be positive", helpers.ErrValidation)
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// A conditional UPDATE locks the product row and rechecks available stock
		// after concurrent updates. The order and stock change commit together.
		result := tx.Model(&product.Product{}).
			Where("id = ? AND remaining_quantity >= ?", order.ProductID, order.Quantity).
			Update("remaining_quantity", gorm.Expr("remaining_quantity - ?", order.Quantity))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			var existing product.Product
			if err := tx.Select("id").First(&existing, order.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("%w: product does not exist", helpers.ErrNotFound)
				}
				return err
			}
			return fmt.Errorf("%w: insufficient product stock", helpers.ErrConflict)
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// Read the price while the stock update still holds the product row lock.
		return tx.Select("id", "price").First(&order.Product, order.ProductID).Error
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return nil, fmt.Errorf("%w: product does not exist", helpers.ErrNotFound)
	}

	return order, err
}
