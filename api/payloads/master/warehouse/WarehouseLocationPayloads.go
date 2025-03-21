package masterwarehousepayloads

import (
	"encoding/json"
)

type SaveWarehouseLocationRequest struct {
	IsActive                      *bool   `json:"is_active"`
	CompanyId                     int     `json:"company_id"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
	WarehouseLocationCode         string  `json:"warehouse_location_code"`
	WarehouseLocationName         string  `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string  `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int     `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64 `json:"warehouse_location_capacity_in_m3"`
}

type WarehouseLocationTemplate struct {
	WarehouseCode         string `json:"warehouse_code"`
	WarehouseLocationCode string `json:"warehouse_location_code"`
	WarehouseLocationName string `json:"warehouse_location_name"`
	Validation            string `json:"validation"`
}

type ProcessWarehouseLocationTemplate struct {
	Data []WarehouseLocationTemplate `json:"data"`
}

type UpdateWarehouseLocationRequest struct {
	IsActive                      *bool   `json:"is_active"`
	WarehouseLocationId           int     `json:"warehouse_location_id"`
	CompanyId                     int     `json:"company_id"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
	WarehouseLocationCode         string  `json:"warehouse_location_code"`
	WarehouseLocationName         string  `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string  `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int     `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64 `json:"warehouse_location_capacity_in_m3"`
}

type GetItemGroupLoction struct {
	WarehouseGroupId int    `json:"warehouse_group_id"`
	WarehouseName    string `json:"warehouse_name"`
	WarehouseCode    string `json:"warehouse_code"`
}

type GetWarehouseLocationRequest struct {
	IsActive                      bool    `json:"is_active"`
	WarehouseLocationId           int     `json:"warehouse_location_id"`
	CompanyId                     int     `json:"company_id"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
	WarehouseLocationCode         string  `json:"warehouse_location_code"`
	WarehouseLocationName         string  `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string  `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int     `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64 `json:"warehouse_location_capacity_in_m3"`
}

type GetWarehouseLocationResponse struct {
	IsActive                      bool    `json:"is_active"`
	WarehouseLocationId           int     `json:"warehouse_location_id"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
	WarehouseLocationCode         string  `json:"warehouse_location_code"`
	WarehouseLocationName         string  `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string  `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int     `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64 `json:"warehouse_location_capacity_in_m3"`
}

type GetWarehouseLocationPreviewResponse struct {
	WarehouseCode         string `json:"warehouse_code"`
	WarehouseLocationCode string `json:"warehouse_location_code"`
	WarehouseLocationName string `json:"warehouse_location_name"`
	Validation            string `json:"validation"`
}

func ConvertWarehouseLocationMapToStruct(maps []map[string]any) ([]GetWarehouseLocationResponse, error) {
	var result []GetWarehouseLocationResponse
	// Marshal the maps into JSON
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type GetAllWarehouseLocationRequest struct {
	IsActive                      string `json:"is_active"`
	WarehouseLocationId           string `json:"warehouse_location_id"`
	CompanyId                     string `json:"company_id"`
	WarehouseGroupId              string `json:"warehouse_group_id"`
	WarehouseLocationCode         string `json:"warehouse_location_code"`
	WarehouseLocationName         string `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence string `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 string `json:"warehouse_location_capacity_in_m3"`
}

type GetAllWarehouseLocationResponse struct {
	IsActive                      bool    `json:"is_active"`
	WarehouseLocationId           int     `json:"warehouse_location_id"`
	CompanyId                     int     `json:"company_id"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
	WarehouseGroupName            string  `json:"warehouse_group_name"`
	WarehouseId                   int     `json:"warehouse_id"`
	WarehouseCode                 string  `json:"warehouse_code"`
	WarehouseName                 string  `json:"warehouse_name"`
	WarehouseLocationCode         string  `json:"warehouse_location_code"`
	WarehouseLocationName         string  `json:"warehouse_location_name"`
	WarehouseLocationDetailName   string  `json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int     `json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64 `json:"warehouse_location_capacity_in_m3"`
}
