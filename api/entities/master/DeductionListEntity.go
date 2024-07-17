package masterentities

import "time"

var CreateDeductionTable = "mtr_deduction"

type DeductionList struct {
	IsActive        bool              `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DeductionId     int               `gorm:"column:deduction_id;size:30;not null;primaryKey;autoincrement" json:"deduction_id"`
	DeductionName   string            `gorm:"column:deduction_name;size:10;not null" json:"deduction_name"`
	DeductionCode   string            `gorm:"column:deduction_code;size:50;not null" json:"deduction_code"`
	EffectiveDate   time.Time         `gorm:"column:effective_date;not null" json:"effective_date"`
	DeductionDetail []DeductionDetail `gorm:"foreignkey:DeductionId;referncces:DeductionId"`
}

func (*DeductionList) TableName() string {
	return CreateDeductionTable
}
