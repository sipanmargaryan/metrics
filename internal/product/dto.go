package product

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,notblank" example:"Office supplies"`
	Quantity    int64  `json:"quantity" binding:"gte=0" example:"100"`
	WarehouseID uint   `json:"warehouse_id" binding:"required,gt=0" example:"1"`
}

type ProductResponse struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Office supplies"`
	Quantity    int64  `json:"quantity" example:"100"`
	WarehouseID uint   `json:"warehouse_id" example:"1"`
}
