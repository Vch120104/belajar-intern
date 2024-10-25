package masteritementities

var CreateItemLevel4Table = "mtr_item_level_4"

type ItemLevel4 struct {
	IsActive       bool `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemLevel4Id   int  `gorm:"column:item_level_4_id;size:30;not null;primaryKey" json:"item_level_4_id"`
	ItemLevel3Id   int  `gorm:"column:item_level_3_id;size:30;not null" json:"item_level_3_id"`
	ItemLevel3     ItemLevel3
	ItemLevel4Code string `gorm:"column:item_level_4_code;not null;size:10;uniqueindex:idx_item_level_4" json:"item_level_4_code"`
	ItemLevel4Name string `gorm:"column:item_level_4_name;not null;size:100;uniqueindex:idx_item_level_4" json:"item_level_4_name"`
}

func (*ItemLevel4) TableName() string {
	return CreateItemLevel4Table
}
