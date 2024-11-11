package masteritementities

var CreateUomItemTable = "mtr_uom_item"

type UomItem struct {
	IsActive          bool    `gorm:"column:is_active;not null"        json:"is_active"`
	UomItemId         int     `gorm:"column:uom_item_id;not null;primaryKey;size:30"        json:"uom_item_id"`
	ItemId            int     `gorm:"column:item_id;not null;size:30;uniqueindex:idx_uom_item"        json:"item_id"`
	UomSourceTypeCode string  `gorm:"column:uom_source_type_code;not null"        json:"uom_source_type_code"`
	UomTypeId         int     `gorm:"column:uom_type_id;size:30;not null" json:"uom_type_id"`
	UomTypeCode       string  `gorm:"column:uom_type_code"        json:"uom_type_code"`
	SourceUomId       int     `gorm:"column:source_uom_id;not null;size:30"        json:"source_uom_id"`
	TargetUomId       int     `gorm:"column:target_uom_id;not null;size:30"        json:"target_uom_id"`
	SourceConvertion  float64 `gorm:"column:source_convertion;not null"        json:"source_convertion"`
	TargetConvertion  float64 `gorm:"column:target_convertion;not null"        json:"target_convertion"`
	UomType           UomType `gorm:"foreignKey:UomTypeId;references:UomTypeId"`
	SourceUom         Uom     `gorm:"foreignKey:SourceUomId;references:UomId"`
	TargetUom         Uom     `gorm:"foreignKey:TargetUomId;references:UomId"`
}

func (*UomItem) TableName() string {
	return CreateUomItemTable
}
