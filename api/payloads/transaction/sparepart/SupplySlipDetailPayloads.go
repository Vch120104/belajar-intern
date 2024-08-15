package transactionsparepartpayloads

type SupplySlipDetailResponse struct {
	SupplyDetailSystemNumber int32   `json:"supply_detail_system_number"`
	OperationCode            string  `json:"operation_code"`
	OperationName            string  `json:"operation_name"`
	ItemCode                 string  `json:"item_code"`
	ItemName                 string  `json:"item_name"`
	WarehouseGroupCode       string  `json:"warehouse_group_code"`
	WarehouseCode            string  `json:"warehouse_code"`
	WarehouseLocationCode    string  `json:"warehouse_location_code"`
	UomCode                  string  `json:"uom_code"`
	QuantitySupply           float64 `json:"quantity_supply"`
	QuantityDemand           float64 `json:"quantity_demand"`
	PR                       bool    `json:"pr"`
	QuantityPR               float32 `json:"quantity_pr"`
}

type SupplySlipDetailByHeaderIdResponse struct {
	SupplyDetailSystemNumber int32   `json:"supply_detail_system_number"`
	OperationCode            string  `json:"operation_code"`
	ItemCode                 string  `json:"item_code"`
	ItemName                 string  `json:"item_name"`
	UomCode                  string  `json:"uom_code"`
	WarehouseGroupCode       string  `json:"warehouse_group_code"`
	WarehouseCode            string  `json:"warehouse_code"`
	WarehouseLocationCode    string  `json:"warehouse_location_code"`
	QuantitySupply           float32 `json:"quantity_supply"`
	QuantityPR               float32 `json:"quantity_pr"`
}


type SupplySlipDetailsResponse struct {
	Page       int                                  `json:"page"`
	Limit      int                                  `json:"limit"`
	TotalPages int                                  `json:"total_pages"`
	TotalRows  int                                  `json:"total_rows"`
	Data       []SupplySlipDetailByHeaderIdResponse `json:"data"`
}
