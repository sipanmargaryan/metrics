package warehouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"metrics/internal/helpers"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *Warehouse) (*Warehouse, error)
	List(ctx context.Context) ([]Warehouse, error)
}

type warehouseRepository struct {
	db *gorm.DB
}

func InitWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (repo *warehouseRepository) Create(ctx context.Context, warehouse *Warehouse) (*Warehouse, error) {
	err := repo.db.WithContext(ctx).Create(warehouse).Error
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, fmt.Errorf("%w: warehouse name already exists", helpers.ErrConflict)
	}
	return warehouse, err
}

func (repo *warehouseRepository) List(ctx context.Context) ([]Warehouse, error) {
	warehouses := make([]Warehouse, 0)
	err := repo.db.WithContext(ctx).Order("id").Find(&warehouses).Error
	return warehouses, err
}
