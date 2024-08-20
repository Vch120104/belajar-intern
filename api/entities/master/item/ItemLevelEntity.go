package masteritementities

var CreateItemLevelTable = "mtr_item_level"

type ItemLevel struct {
	IsActive        bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemLevelId     int    `gorm:"column:item_level_id;size:30;not null;primaryKey" json:"item_level_id"`
	ItemLevel       string `gorm:"column:item_level;not null;type:varchar(1);uniqueindex:idx_item_level" json:"item_level"`
	ItemClassId     int    `gorm:"column:item_class_id;size:30;not null"        json:"item_class_id"`
	ItemLevelParent int    `gorm:"column:item_level_parent;not null;size:30" json:"item_level_parent"`
	ItemLevelCode   string `gorm:"column:item_level_code;not null;type:varchar(10);uniqueindex:idx_item_level" json:"item_level_code"`
	ItemLevelName   string `gorm:"column:item_level_name;not null;type:varchar(100);uniqueindex:idx_item_level" json:"item_level_name"`
}

func (*ItemLevel) TableName() string {
	return CreateItemLevelTable
}
