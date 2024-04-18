package masteritempayloads

type ItemLocationResponse struct {
	ItemLocationId     int    `json:"item_location_id"`
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	ItemId             int    `json:"item_id"`
	ItemCode           string `json:"item_code"`
	ItemName           string `json:"item_name"`
}

type ItemLocWarehouseGroupResponse struct {
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
}

type ItemLocResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}

type ItemLocSourceResponse struct {
	ItemLocationSourceId   int    `json:"item_location_source_id"`
	ItemLocationSourceCode string `json:"item_location_source_code"`
	ItemLocationSourceName string `json:"item_location_source_name"`
}

type ItemLocSourceRequest struct {
	ItemLocationSourceId   int    `json:"item_location_source_id" parent_entity:"mtr_item_location_source" main_table:"mtr_item_location_source"`
	ItemLocationSourceCode string `json:"item_location_source_code" parent_entity:"mtr_item_location_source"`
	ItemLocationSourceName string `json:"item_location_source_name" parent_entity:"mtr_item_location_source"`
}

type ItemLocationRequest struct {
	ItemLocationId   int `json:"item_location_id" parent_entity:"mtr_item_location" main_table:"mtr_item_location"`
	WarehouseGroupId int `json:"warehouse_group_id" parent_entity:"mtr_item_location"`
	ItemId           int `json:"item_id" parent_entity:"mtr_item_location"`
}

type ItemLocationDetailResponse struct {
	ItemLocationDetailId int `json:"item_location_detail_id"`
	ItemLocationId       int `json:"item_location_id"`
	ItemId               int `json:"item_id"`
	ItemLocationSourceId int `json:"item_location_source_id"`
}

type ItemLocationDetailRequest struct {
	ItemLocationDetailId int `json:"item_location_detail_id" parent_entity:"mtr_item_location_detail" main_table:"mtr_item_location_detail"`
	ItemLocationId       int `json:"item_location_id" parent_entity:"mtr_item_location_detail"`
	ItemId               int `json:"item_id" parent_entity:"mtr_item_location_detail"`
	ItemLocationSourceId int `json:"item_location_source_id" parent_entity:"mtr_item_location_detail"`
}
