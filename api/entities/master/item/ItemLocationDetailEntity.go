package masteritementities

var CreateItemLocationDetailTable = "mtr_item_location_detail"

type ItemLocationDetail struct {
	ItemLocationDetailId int    `gorm:"column:item_location_detail_id;size:30;primaryKey" json:"item_location_detail_id"`
	ItemLocationId       int    `gorm:"column:item_location_id;size:30;not null" json:"item_location_id"`
	ItemId               int    `gorm:"column:item_id;size:30;not null" json:"item_id"`
	ItemCode             string `gorm:"column:item_code;not null;type:varchar(100)" json:"item_code"`
	ItemName             string `gorm:"column:item_name;not null;type:varchar(100)" json:"item_name"`
	LocationCode         string `gorm:"column:item_location_detail_code;not null;type:varchar(100)" json:"item_location_detail_code"`
	LocationName         string `gorm:"column:item_location_detail_name;not null;type:varchar(100)" json:"item_location_detail_name"`
}

func (*ItemLocationDetail) TableName() string {
	return CreateItemLocationDetailTable
}
