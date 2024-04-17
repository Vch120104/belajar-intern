package masteritementities

var (
	CreateItemLocationTable = "mtr_item_location"
)

type ItemLocation struct {
	ItemLocationId     int    `gorm:"column:item_location_id;size:30;primaryKey" json:"item_location_id"`
	WarehouseGroupId   int    `gorm:"column:warehouse_group_id;size:30;not null" json:"warehouse_group_id"`
	WarehouseGroupCode string `gorm:"column:warehouse_group_code;not null;type:varchar(100)" json:"warehouse_group_code"`
	ItemId             int    `gorm:"column:item_id;size:30;not null" json:"item_id"`
	ItemCode           string `gorm:"column:item_code;not null;type:varchar(100)" json:"item_code"`
	ItemName           string `gorm:"column:item_name;not null;type:varchar(100)" json:"item_name"`
}

func (*ItemLocation) TableName() string {
	return CreateItemLocationTable
}
