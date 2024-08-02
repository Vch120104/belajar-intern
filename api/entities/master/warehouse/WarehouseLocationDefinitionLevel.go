package masterwarehouseentities

const TableNameWarehouseLocationDefinitionLevel = "mtr_warehouse_location_definition_level"

type WarehouseLocationDefinitionLevel struct {
	WarehouseLocationDefinitionLevelId          int    `gorm:"column:warehouse_location_definition_level_id;size:30;not null;primaryKey" json:"warehouse_location_definition_level_id"`
	WarehouseLocationDefinitionLevelDescription string `gorm:"column:warehouse_location_definition_level_description;not null;type:varchar(100)" json:"warehouse_location_definition_level_description"`
}

func (*WarehouseLocationDefinitionLevel) TableName() string {
	return TableNameWarehouseLocationDefinitionLevel
}
