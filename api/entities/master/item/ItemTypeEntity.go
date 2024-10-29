package masteritementities

var CreateItemTypeTable = "mtr_item_type"

type ItemType struct {
	IsActive     bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemTypeId   int    `gorm:"column:item_type_id;size:30;not null;primaryKey" json:"item_type_id"`
	ItemTypeCode string `gorm:"column:item_type_code;not null;type:varchar(10);" json:"item_type_code"`
	ItemTypeName string `gorm:"column:item_type_name;not null;type:varchar(100);" json:"item_type_name"`
}

func (*ItemType) TableName() string {
	return CreateItemTypeTable
}
