package product

import (
	"time"

	"metrics/internal/warehouse"
)

type Product struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string              `gorm:"not null;uniqueIndex:products_warehouse_name_key"`
	Quantity    int64               `gorm:"not null;default:0;check:products_quantity_nonnegative,quantity >= 0"`
	WarehouseID uint                `gorm:"not null;index;uniqueIndex:products_warehouse_name_key"`
	Warehouse   warehouse.Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnDelete:RESTRICT"`
}
