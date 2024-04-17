package masteritempayloads

type ItemLocationRequest struct {
	ItemLocationId     int    `json:"item_location_id"`
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	ItemId             int    `json:"item_id"`
	ItemCode           string `json:"item_code"`
	ItemName           string `json:"item_name"`
}

type ItemLocWarehouseGroupResponse struct {
	WarehouseGroupId   string `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
}

type ItemLocResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

type ItemLocationResponse struct {
	ItemLocationId     int    `json:"item_location_id" parent_entity:"mtr_item_location" main_table:"mtr_item_location"`
	WarehouseGroupId   int    `json:"warehouse_group_id" parent_entity:"mtr_item_location"`
	WarehouseGroupCode string `json:"warehouse_group_code" parent_entity:"mtr_item_location"`
	ItemId             int    `json:"item_id" parent_entity:"mtr_item_location"`
	ItemCode           string `json:"item_code" parent_entity:"mtr_item_location"`
	ItemName           string `json:"item_name" parent_entity:"mtr_item_location"`
}
