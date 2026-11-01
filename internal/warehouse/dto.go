package warehouse

type CreateWarehouseRequest struct {
	Name string `json:"name" binding:"required,notblank" example:"Yerevan warehouse"`
}

type WarehouseResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Yerevan warehouse"`
}
