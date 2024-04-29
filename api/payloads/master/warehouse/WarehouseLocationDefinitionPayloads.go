package masterwarehousepayloads

type WarehouseLocationDefinitionRequest struct {
	IsActive                               bool   `json:"is_active" parent_entity:"mtr_warehouse_location_definition" main_table:"mtr_warehouse_location_definition"`
	WarehouseLocationDefinitionId          int    `json:"warehouse_location_definition_id" parent_entity:"mtr_warehouse_location_definition"`
	WarehouseLocationDefinitionLevelId     int    `json:"warehouse_location_definition_level_id" parent_entity:"mtr_warehouse_location_definition"`
	WarehouseLocationDefinitionLevelCode   string `json:"warehouse_location_definition_level_code" parent_entity:"mtr_warehouse_location_definition"`
	WarehouseLocationDefinitionDescription string `json:"warehouse_location_definition_description" parent_entity:"mtr_warehouse_location_definition"`
}

type WarehouseLocationDefinitionResponse struct {
	IsActive                               bool   `json:"is_active"`
	WarehouseLocationDefinitionId          int    `json:"warehouse_location_definition_id"`
	WarehouseLocationDefinitionLevelId     int    `json:"warehouse_location_definition_level_id"`
	WarehouseLocationDefinitionLevelCode   string `json:"warehouse_location_definition_level_code"`
	WarehouseLocationDefinitionDescription string `json:"warehouse_location_definition_description"`
}

type WarehouseLocationDefinitionLevelResponse struct {
	WarehouseLocationDefinitionLevelId          int    `json:"warehouse_location_definition_level_id"`
	WarehouseLocationDefinitionLevelDescription string `json:"warehouse_location_definition_level_description"`
}

type WarehouseLocationDefinitionLevelRequest struct {
	WarehouseLocationDefinitionLevelId          int    `json:"warehouse_location_definition_level_id" parent_entity:"mtr_warehouse_location_definition_level" main_table:"mtr_warehouse_location_definition_level"`
	WarehouseLocationDefinitionLevelDescription string `json:"warehouse_location_definition_level_description" parent_entity:"mtr_warehouse_location_definition_level"`
}
