package order

import "time"

// OrderCreatedEvent contains the order details published to Kafka.
// Price is the unit price in currency minor units.
type OrderCreatedEvent struct {
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Price     int64     `json:"price"`
	OrderType string    `json:"order_type"`
	CreatedAt time.Time `json:"created_at"`
}
