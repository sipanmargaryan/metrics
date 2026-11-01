package product

func toProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		WarehouseID: product.WarehouseID,
	}
}
