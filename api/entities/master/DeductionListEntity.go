package masterentities

import "time"

var CreateDeductionTable = "mtr_deduction_list"

type DeductionList struct {
	IsActive          bool      `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DeductionListId   int       `gorm:"column:deduction_list_id;size:30;not null;primaryKey" json:"deduction_list_id"`
	DeductionName     string    `gorm:"column:deduction_name;size:10;not null" json:"deduction_name"`
	EffectiveDate     time.Time `gorm:"column:effective_date;not null" json:"effective_date"`
}

func (*DeductionList) TableName() string {
	return CreateDeductionTable
}
