package order

type CreateOrderRequest struct {
	ProductID uint  `json:"product_id" binding:"required,gt=0" example:"1"`
	Quantity  int64 `json:"quantity" binding:"required,gt=0" example:"10"`
}

type OrderCreatedResponse struct {
	ID uint `json:"id" example:"1"`
}
