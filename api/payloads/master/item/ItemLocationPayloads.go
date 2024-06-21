package masteritempayloads

type ItemLocationRequest struct {
	ItemLocationId   int `json:"item_location_id" parent_entity:"mtr_item_location" main_table:"mtr_item_location"`
	WarehouseGroupId int `json:"warehouse_group_id" parent_entity:"mtr_item_location"`
	WarehouseId      int `json:"warehouse_id" parent_entity:"mtr_item_location"`
	ItemId           int `json:"item_id" parent_entity:"mtr_item_location"`
}

type ItemLocationResponse struct {
	ItemLocationId     int    `json:"item_location_id"`
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseId        int    `json:"warehouse_id"`
	WarehouseCode      string `json:"warehouse_code"`
	WarehouseName      string `json:"warehouse_name"`
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

type ItemLocationWarehouseResponse struct {
	WarehouseId   int    `json:"warehouse_id"`
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
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

type ItemLocationGetAllResponse struct {
	ItemLocationId        int    `json:"item_location_id" parent_entity:"mtr_location_item" main_table:"mtr_location_item"`
	ItemId                int    `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode              string `json:"item_code" parent_entity:"mtr_item"`
	ItemName              string `json:"item_name" parent_entity:"mtr_item"`
	StockOpname           bool   `json:"stock_opname" parent_entity:"mtr_location_item"`
	WarehouseId           int    `json:"warehouse_id" parent_entity:"mtr_warehouse_master" references:"mtr_warehouse_master"`
	WarehouseName         string `json:"warehouse_name" parent_entity:"mtr_warehouse_master"`
	WarehouseCode         string `json:"warehouse_code" parent_entity:"mtr_warehouse_master"`
	WarehouseGroupId      int    `json:"warehouse_group_id" parent_entity:"mtr_warehouse_group" references:"mtr_warehouse_group"`
	WarehouseGroupName    string `json:"warehouse_group_name" parent_entity:"mtr_warehouse_group"`
	WarehouseGroupCode    string `json:"warehouse_group_code" parent_entity:"mtr_warehouse_group"`
	WarehouseLocationId   int    `json:"warehouse_location_id" parent_entity:"mtr_warehouse_location" references:"mtr_warehouse_location"`
	WarehouseLocationName string `json:"warehouse_location_name" parent_entity:"mtr_warehouse_location"`
	WarehouseLocationCode string `json:"warehouse_location_code" parent_entity:"mtr_warehouse_location"`
}

type ItemLocationGetByIdResponse struct {
	ItemLocationId        int    `json:"item_location_id" parent_entity:"mtr_item_location" main_table:"mtr_item_location"`
	WarehouseGroupId      int    `json:"warehouse_group_id" parent_entity:"mtr_item_location"`
	WarehouseId           int    `json:"warehouse_id" parent_entity:"mtr_item_location"`
	ItemId                int    `json:"item_id" parent_entity:"mtr_item_location"`
	ItemCode              string `json:"item_code" parent_entity:"mtr_item"`
	ItemName              string `json:"item_name" parent_entity:"mtr_item"`
	WarehouseLocationId   int    `json:"warehouse_location_id" parent_entity:"mtr_item_location"`
	WarehouseLocationName string `json:"warehouse_location_name" parent_entity:"mtr_warehouse_location"`
	WarehouseLocationCode string `json:"warehouse_location_code" parent_entity:"mtr_warehouse_location"`
}

type SaveItemlocation struct {
	ItemLocationId      int `json:"item_location_id"`
	WarehouseGroupId    int `json:"warehouse_group_id"`
	ItemId              int `json:"item_id"`
	WarehouseId         int `json:"warehouse_id"`
	WarehouseLocationId int `json:"warehouse_location_id"`
}
