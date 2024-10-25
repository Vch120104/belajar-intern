package masteritementities

var CreateItemLevel3Table = "mtr_item_level_3"

type ItemLevel3 struct {
	IsActive       bool `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemLevel3Id   int  `gorm:"column:item_level_3_id;size:30;not null;primaryKey" json:"item_level_3_id"`
	ItemLevel2Id   int  `gorm:"column:item_level_2_id;size:30;not null" json:"item_level_2_id"`
	ItemLevel2     ItemLevel2
	ItemLevel3Code string `gorm:"column:item_level_3_code;not null;size:10;uniqueindex:idx_item_level_3" json:"item_level_3_code"`
	ItemLevel3Name string `gorm:"column:item_level_3_name;not null;size:100;uniqueindex:idx_item_level_3" json:"item_level_3_name"`
}

func (*ItemLevel3) TableName() string {
	return CreateItemLevel3Table
}
