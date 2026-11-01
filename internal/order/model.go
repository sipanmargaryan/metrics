package order

import (
	"time"

	"metrics/internal/product"
)

type Order struct {
	ID uint `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Quantity  int64           `gorm:"not null;default:1"`
	ProductID uint            `gorm:"not null;index"`
	Product   product.Product `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT"`
}
