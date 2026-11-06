package product

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,notblank" example:"Office supplies"`
	Price       int64  `json:"price" binding:"gte=0" example:"1999"` // Currency minor units.
	Quantity    int64  `json:"quantity" binding:"gte=0" example:"100"`
	WarehouseID uint   `json:"warehouse_id" binding:"required,gt=0" example:"1"`
}

type ProductResponse struct {
	ID                uint   `json:"id" example:"1"`
	Name              string `json:"name" example:"Office supplies"`
	Price             int64  `json:"price" example:"1999"` // Currency minor units.
	InitialQuantity   int64  `json:"initial_quantity" example:"100"`
	RemainingQuantity int64  `json:"remaining_quantity" example:"90"`
	WarehouseID       uint   `json:"warehouse_id" example:"1"`
}
