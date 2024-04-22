package masterentities

var CreateDeductionDetailTable = "mtr_deduction_detail"

type DeductionDetail struct {
	IsActive             bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DeductionDetailId    int     `gorm:"column:deduction_detail_id;size:30;not null;primaryKey" json:"deduction_detail_id"`
	DeductionListId      int     `gorm:"column:deduction_list_id;size:30;not null" json:"deduction_list_id"`
	DeductionList        DeductionList 
	DeductionDetailCode  string  `gorm:"column:deduction_detail_code;size:50;not null" json:"deduction_detail_code"`
	DeductionDetailLevel int     `gorm:"column:deduction_detail_level;size:30;not null" json:"deduction_detail_level"`
	DeductionPercent     float64 `gorm:"column:deduction_percent;not null" json:"deduction_percent"`
}

func (*DeductionDetail) TableName() string {
	return CreateDeductionDetailTable
}
