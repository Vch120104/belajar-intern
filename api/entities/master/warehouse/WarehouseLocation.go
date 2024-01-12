package masterwarehouseentities

const TableNameWarehouseLocation = "mtr_warehouse_location"

type WarehouseLocation struct {
	IsActive                      *bool          `gorm:"column:is_active;default:true;not null" json:"is_active"`
	WarehouseLocationId           int            `gorm:"column:warehouse_location_id;not null;primaryKey" json:"warehouse_location_id"`
	CompanyId                     int            `gorm:"column:company_id;not null" json:"company_id"`
	WarehouseGroupId              int            `gorm:"column:warehouse_group_id;not null" json:"warehouse_group_id"`
	WarehouseLocationCode         string         `gorm:"column:warehouse_location_code;not null;type:varchar(5)" json:"warehouse_location_code"`
	WarehouseLocationName         string         `gorm:"column:warehouse_location_name;not null;type:varchar(100)" json:"warehouse_location_name"`
	WarehouseLocationDetailName   string         `gorm:"column:warehouse_location_detail_name;not null;type:varchar(100)" json:"warehouse_location_detail_name"`
	WarehouseLocationPickSequence int            `gorm:"column:warehouse_location_pick_sequence;not null" json:"warehouse_location_pick_sequence"`
	WarehouseLocationCapacityInM3 float64        `gorm:"column:warehouse_location_capacity_in_m3;not null" json:"warehouse_location_capacity_in_m3"`
}

func (*WarehouseLocation) TableName() string {
	return TableNameWarehouseLocation
}
