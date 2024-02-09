package masteritementities

import "time"

const CreateItemSubstitute = "mtr_item_substitute"

type ItemSubstitute struct {
	IsActive           bool      `gorm:"column:is_active;not null;default:True" json:"is_active"`
	SubstituteTypeCode string    `gorm:"column:substitute_type_code;not null;unique;size:2" json:"substitute_type_code"`
	ItemSubstituteId   int       `gorm:"column:item_substitute_id;not null;size:30;primaryKey" json:"item_substitute_is"`
	EffectiveDate      time.Time `gorm:"column:effective_date;not null" json:"effective_date"`
	ItemId             int       `gorm:"column:item_id;not null;unique" json:"item_id"`
}

func (*ItemSubstitute) TableName() string {
	return CreateItemSubstitute
}
