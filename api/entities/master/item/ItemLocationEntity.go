package masteritementities

var (
	CreateItemLocationTable = "mtr_location_item"
)

type ItemLocation struct {
	ItemLocationId      int  `gorm:"column:item_location_id;size:30;primaryKey" json:"item_location_id"`
	WarehouseGroupId    int  `gorm:"column:warehouse_group_id;size:30;not null" json:"warehouse_group_id"`
	ItemId              int  `gorm:"column:item_id;size:30;not null" json:"item_id"`
	WarehouseId         int  `gorm:"column:warehouse_id;size:30;not null " json:"warehouse_id"`
	WarehouseLocationId int  `gorm:"column:warehouse_location_id;size:30;not null" json:"warehouse_location_id"`
	StockOpname         bool `gorm:"column:stock_opname;default:true" json:"stock_opname"`
}

func (*ItemLocation) TableName() string {
	return CreateItemLocationTable
}
