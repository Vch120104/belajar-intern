package masteritementities

var CreateItemLevel1Table = "mtr_item_level_1"

type ItemLevel1 struct {
	IsActive       bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemLevel1Id   int    `gorm:"column:item_level_1_id;size:30;not null;primaryKey" json:"item_level_1_id"`
	ItemLevel1Code string `gorm:"column:item_level_1_code;not null;size:10;uniqueindex:idx_item_level_1" json:"item_level_1_code"`
	ItemLevel1Name string `gorm:"column:item_level_1_name;not null;size:100" json:"item_level_1_name"`
	ItemClassId    int    `gorm:"column:item_class_id;size:30;not null;uniqueindex:idx_item_level_1" json:"item_class_id"`
}

func (*ItemLevel1) TableName() string {
	return CreateItemLevel1Table
}
