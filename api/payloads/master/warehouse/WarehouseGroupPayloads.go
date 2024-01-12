package masterwarehousepayloads

type SaveWarehouseGroupRequest struct {
	IsActive           bool   `json:"is_active"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
	ProfitCenterId     int    `json:"profit_center_id"`
}

type UpdateWarehouseGroupRequest struct {
	IsActive           bool   `json:"is_active"`
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
	ProfitCenterId     int    `json:"profit_center_id"`
}

type GetWarehouseGroupResponse struct {
	IsActive           bool   `json:"is_active"`
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
	ProfitCenterId     int    `json:"profit_center_id"`
}

type GetAllWarehouseGroupRequest struct {
	IsActive           string `json:"is_active"`
	WarehouseGroupId   string `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
	ProfitCenterId     string `json:"profit_center_id"`
}

type ProfitCenterResponse struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"` 
	ProfitCenterName string `json:"profit_center_name"` 
}