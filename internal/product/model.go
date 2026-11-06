package product

import (
	"time"

	"metrics/internal/warehouse"
)

type Product struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;uniqueIndex:products_warehouse_name_key"`
	// Price is stored in the currency's minor units (for example, cents).
	Price             int64               `gorm:"not null;default:0;check:products_price_nonnegative,price >= 0"`
	InitialQuantity   int64               `gorm:"not null;default:0;check:products_initial_quantity_nonnegative,initial_quantity >= 0"`
	RemainingQuantity int64               `gorm:"not null;default:0;check:products_remaining_quantity_valid,remaining_quantity >= 0 AND remaining_quantity <= initial_quantity"`
	WarehouseID       uint                `gorm:"not null;index;uniqueIndex:products_warehouse_name_key"`
	Warehouse         warehouse.Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnDelete:RESTRICT"`
}
