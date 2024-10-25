package masteritementities

var CreateItemLevel2Table = "mtr_item_level_2"

type ItemLevel2 struct {
	IsActive       bool `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemLevel2Id   int  `gorm:"column:item_level_2_id;size:30;not null;primaryKey" json:"item_level_2_id"`
	ItemLevel1Id   int  `gorm:"column:item_level_1_id;size:30;not null" json:"item_level_1_id"`
	ItemLevel1     ItemLevel1
	ItemLevel2Code string `gorm:"column:item_level_2_code;not null;size:10;uniqueindex:idx_item_level_2" json:"item_level_2_code"`
	ItemLevel2Name string `gorm:"column:item_level_2_name;not null;size:100;uniqueindex:idx_item_level_2" json:"item_level_2_name"`
}

func (*ItemLevel2) TableName() string {
	return CreateItemLevel2Table
}
