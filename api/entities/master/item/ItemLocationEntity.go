package masteritementities

var (
	CreateItemLocationTable = "mtr_item_location"
)

type ItemLocation struct {
	ItemLocationId     int                  `gorm:"column:item_location_id;size:30;primaryKey" json:"item_location_id"`
	WarehouseGroupId   int                  `gorm:"column:warehouse_group_id;size:30;not null" json:"warehouse_group_id"`
	ItemId             int                  `gorm:"column:item_id;size:30;not null" json:"item_id"`
	WarehouseId        int                  `gorm:"column:warehouse_id;size:30;not null " json:"warehouse_id"`
	ItemLocationDetail []ItemLocationDetail `gorm:"foreignKey:ItemLocationId;" json:"location_detail"`
}

func (*ItemLocation) TableName() string {
	return CreateItemLocationTable
}
