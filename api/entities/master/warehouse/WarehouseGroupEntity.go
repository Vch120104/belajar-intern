package masterwarehouseentities

const TableNameWarehouseGroup = "mtr_warehouse_group"

type WarehouseGroup struct {
	IsActive           *bool             `gorm:"column:is_active;default:true;not null" json:"is_active"`
	WarehouseGroupId   int               `gorm:"column:warehouse_group_id;not null;primaryKey;size:30" json:"warehouse_group_id"`
	WarehouseGroupCode string            `gorm:"column:warehouse_group_code;not null;type:varchar(5)" json:"warehouse_group_code"`
	WarehouseGroupName string            `gorm:"column:warehouse_group_name;not null;type:varchar(100)" json:"warehouse_group_name"`
	ProfitCenterId     int               `gorm:"column:profit_center_id;not null;size:30" json:"profit_center_id"`
	WarehouseMaster    WarehouseMaster   `gorm:"foreignKey:WarehouseGroupId;references:WarehouseGroupId" `
	WarehouseLocation  WarehouseLocation `gorm:"foreignKey:WarehouseGroupId;references:warehouse_group_id"`
}

func (*WarehouseGroup) TableName() string {
	return TableNameWarehouseGroup
}
