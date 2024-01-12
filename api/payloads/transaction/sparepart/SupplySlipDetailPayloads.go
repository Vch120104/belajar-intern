package transactionsparepartpayloads

type SupplySlipDetailResponse struct {
	IsActive                          bool    `json:"is_active"`
	SupplySlipDetailSystemNumber      int32   `json:"supply_slip_detail_system_number"`
	SupplySystemNumber                int32   `json:"supply_system_number"` // Update this line
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
