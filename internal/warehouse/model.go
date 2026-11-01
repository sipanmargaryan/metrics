package warehouse

import "time"

type Warehouse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"not null;uniqueIndex" json:"name"`
}
