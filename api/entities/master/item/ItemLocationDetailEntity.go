package masteritementities

var CreateItemLocationDetailTable = "mtr_item_location_detail"

type ItemLocationDetail struct {
	ItemLocationDetailId       int `gorm:"column:item_location_detail_id;size:30;primaryKey" json:"item_location_detail_id"`
	ItemLocationId             int `gorm:"column:item_location_id;size:30;not null" json:"item_location_id"`
	ItemId                     int `gorm:"column:item_id;size:30;not null" json:"item_id"`
	ItemLocationDetailSourceId int `gorm:"column:item_location_source_id;not null;size:30" json:"item_location_source_id"`
}

func (*ItemLocationDetail) TableName() string {
	return CreateItemLocationDetailTable
}
