package masteritementities

const CreateItemSubstituteDetail = "mtr_item_substitute_detail"

type ItemSubstituteDetail struct {
	IsActive               bool `gorm:"column:is_active;not null;default:True" json:"is_active"`
	ItemSubstituteDetailId int  `gorm:"column:item_substitute_detail_id;size:30;not null;primaryKey" json:"item_substitute_detail_id"`
	ItemSubstituteId       int  `gorm:"column:item_substitute_id;size:30;not null" json:"item_substitute_id"`
	ItemSubstitute         ItemSubstitute
	SubstituteItemId       int     `gorm:"column:substitute_item_id;size:30;not null" json:"substitute_item_id"`
	SubstituteItem         Item    `gorm:"foreignKey:SubstituteItemId"`
	Quantity               float64 `gorm:"column:quantity;not null" json:"quantity"`
	Sequence               int     `gorm:"column:sequence;not null;size:30" json:"sequence"`
}

func (*ItemSubstituteDetail) TableName() string {
	return CreateItemSubstituteDetail
}
