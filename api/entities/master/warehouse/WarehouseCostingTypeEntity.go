package masterwarehouseentities

const mtrCostingTypeTableName = "mtr_warehouse_costing_Type"

type WarehouseCostingType struct {
	WarehouseCostingTypeId          int    `gorm:"column:warehouse_costing_type_id;not null;primaryKey" json:"warehouse_costing_type_id"`
	WarehouseCostingTypeCode        string `gorm:"column:warehouse_costing_type_code;not null" json:"warehouse_costing_type_code"`
	WarehouseCostingTypeDescription string `gorm:"column:warehouse_costing_type_description;not null" json:"warehouse_costing_type_description"`
}

func (*WarehouseCostingType) TableName() string {
	return mtrCostingTypeTableName
}
