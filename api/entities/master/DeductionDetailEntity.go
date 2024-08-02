package masterentities

var CreateDeductionDetailTable = "mtr_deduction_detail"

type DeductionDetail struct {
	IsActive             bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DeductionDetailId    int     `gorm:"column:deduction_detail_id;size:30;not null;primaryKey" json:"deduction_detail_id"`
	DeductionId          int     `gorm:"column:deduction_id;size:30;not null" json:"deduction_id"`
	LimitDays            int     `gorm:"column:limit_days;size:30;not null" json:"limit_days"`
	DeductionDetailLevel int     `gorm:"column:deduction_detail_level;size:30;not null" json:"deduction_detail_level"`
	DeductionPercent     float64 `gorm:"column:deduction_percent;not null" json:"deduction_percent"`
}

func (*DeductionDetail) TableName() string {
	return CreateDeductionDetailTable
}
