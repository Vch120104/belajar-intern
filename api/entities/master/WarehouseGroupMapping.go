package masterentities

type WarehouseGroupMappingEntities struct {
	WarehouseGroupMappingId          int    `gorm:"column:warehouse_group_mapping_id;not null;primaryKey" json:"warehouse_group_mapping_id"`
	CompanyId                        int    `gorm:"column:company_id;not null" json:"company_id"`
	WarehouseGroupTypeCode           string `gorm:"column:warehouse_group_type_code;size:50;not null" json:"warehouse_group_type_code"`
	WarehouseGroupId                 int    `gorm:"column:warehouse_group_id;not null" json:"warehouse_group_id"`
	WarehouseGroupMappingDescription string `gorm:"column:warehouse_group_mapping_description;null"        json:"warehouse_group_mapping_description"`
}

func (*WarehouseGroupMappingEntities) TableName() string {
	return "mtr_warehouse_master_mapping"
}
