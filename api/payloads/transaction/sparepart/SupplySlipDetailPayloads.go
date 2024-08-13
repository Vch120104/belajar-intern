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

type SupplySlipDetailRequest struct {
	SupplySystemNumbers               int32   `json:"supply_system_number"`
	SupplySystemLineNumber            int32   `json:"supply_system_line_number"`
	LocationId                        int32   `json:"location_id"`
	WorkOrderItemId                   int32   `json:"work_order_item_id"`
	UnitOutMeasurementId              int32   `json:"unit_out_measurement_id"`
	QuantitySupply                    float64 `json:"quantity_supply"`
	QuantityReturn                    float64 `json:"quantity_return"`
	QuantityDemand                    float64 `json:"quantity_demand"`
	CostOfGoodsSold                   float64 `json:"cost_of_goods_sold"`
	PurchaseRequestSystemNumber       int32   `json:"purchase_request_system_number"`
	PurchaseRequestSystemNumberDetail int32   `json:"purchase_request_system_number_detail"`
	WorkOrderSystemNumber             int32   `json:"work_order_system_number"`
	WorkOrderLineNumberId             int32   `json:"work_order_line_number_id"`
	WarehouseGroupId                  int32   `json:"warehouse_group_id"`
	WarehouseId                       int32   `json:"warehouse_id"`
	QuantityTotal                     int32   `json:"quantity_total"`
}

type SupplySlipDetailsResponse struct {
	Page       int                                  `json:"page"`
	Limit      int                                  `json:"limit"`
	TotalPages int                                  `json:"total_pages"`
	TotalRows  int                                  `json:"total_rows"`
	Data       []SupplySlipDetailByHeaderIdResponse `json:"data"`
}
