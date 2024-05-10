package masterentities

var CreateAgreementDiscountTable = "mtr_agreement_discount_detail"

type AgreementDiscount struct {
	AgreementDiscountId int     `gorm:"column:agreement_discount_id;size:30;not null;primaryKey" json:"agreement_discount_id"`
	AgreementId         int     `gorm:"column:agreement_id;size:30;not null;" json:"agreement_id"`
	LineTypeId          int     `gorm:"column:line_type_id;size:30;not null;" json:"line_type_id"`
	MinValue            int     `gorm:"column:min_value;size:30;not null;" json:"min_value"`
	DiscountPercent     float32 `gorm:"column:discount_percent;size:30;not null;" json:"discount_percent"`
	DiscountRemarks     string  `gorm:"column:discount_remarks;size:50;not null" json:"discount_remarks"`
}

func (*AgreementDiscount) TableName() string {
	return CreateAgreementDiscountTable
}
