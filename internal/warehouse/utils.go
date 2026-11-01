package warehouse

func toWarehouseResponse(warehouse Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:   warehouse.ID,
		Name: warehouse.Name,
	}
}
