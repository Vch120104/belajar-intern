package masteritementities

var (
	CreateItemLocationSourceTable = "mtr_item_location_source"
)

type ItemLocationSource struct {
	ItemLocationSourceId   int    `gorm:"column:item_location_source_id;size:30;primaryKey" json:"item_location_source_id"`
	ItemLocationSourceCode string `gorm:"column:item_location_source_code;size:30;not null" json:"item_location_source_code"`
	ItemLocationSourceName string `gorm:"column:item_location_source_name;size:30;not null" json:"item_location_source_name"`
	IsActive               bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
}

func (*ItemLocationSource) TableName() string {
	return CreateItemLocationSourceTable
}
