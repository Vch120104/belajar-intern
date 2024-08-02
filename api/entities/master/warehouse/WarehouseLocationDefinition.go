package masterwarehouseentities

const TableNameWarehouseLocationDefinition = "mtr_warehouse_location_definition"

type WarehouseLocationDefinition struct {
	IsActive                               bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	WarehouseLocationDefinitionId          int    `gorm:"column:warehouse_location_definition_id;size:30;not null;primaryKey" json:"warehouse_location_definition_id"`
	WarehouseLocationDefinitionLevelId     int    `gorm:"column:warehouse_location_definition_level_id;size:30;not null;" json:"warehouse_location_definition_level_id"`
	WarehouseLocationDefinitionLevelCode   string `gorm:"column:warehouse_location_definition_level_code;not null;type:varchar(100)" json:"warehouse_location_definition_level_code"`
	WarehouseLocationDefinitionDescription string `gorm:"column:warehouse_location_definition_description;not null;type:varchar(100)" json:"warehouse_location_definition_description"`
}

func (*WarehouseLocationDefinition) TableName() string {
	return TableNameWarehouseLocationDefinition
}
