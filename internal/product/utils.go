package product

func toProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:                product.ID,
		Name:              product.Name,
		Price:             product.Price,
		InitialQuantity:   product.InitialQuantity,
		RemainingQuantity: product.RemainingQuantity,
		WarehouseID:       product.WarehouseID,
	}
}
