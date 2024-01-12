package masteritementities

var CreateItemClassTable = "mtr_item_class"

type ItemClass struct {
	IsActive      bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	ItemClassId   int    `gorm:"column:item_class_id;not null;primaryKey"  json:"item_class_id"`
	ItemClassCode string `gorm:"column:item_class_code;index:idx_item_class_code;unique;type:varchar(10)" json:"item_class_code"`
	ItemGroupID   int    `gorm:"column:item_group_id;not null;" json:"item_group_id"` //FK with mtr_item_group common-general service
	LineTypeID    int    `gorm:"column:line_type_id;not null" json:"line_type_id"`    //FK with mtr_line_type common-general service
	ItemClassName string `gorm:"column:item_class_name;not null"  json:"item_class_name"`
}

func (*ItemClass) TableName() string {
	return CreateItemClassTable
}
