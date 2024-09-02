package transactionsparepartpayloads

type SupplySlipReturnDetailResponse struct {
	SupplyReturnDetailSystemNumber int32   `json:"supply_return_detail_system_number"`
	SupplyDetailSystemNumber       int32   `json:"supply_detail_system_number"`
	OperationCode                  string  `json:"operation_code"`
	ItemCode                       string  `json:"item_code"`
	ItemName                       string  `json:"item_name"`
	WarehouseGroupCode             string  `json:"warehouse_group_code"`
	WarehouseCode                  string  `json:"warehouse_code"`
	WarehouseLocationCode          string  `json:"warehouse_location_code"`
	UomCode                        string  `json:"uom_code"`
	QuantitySupply                 float32 `json:"quantity_supply"`
	QuantityReturn                 float32 `json:"quantity_return"`
}
