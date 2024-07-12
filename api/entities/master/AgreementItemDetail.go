package masterentities

var CreateAgreementItemTable = "mtr_agreement_item_detail"

type AgreementItemDetail struct {
	AgreementItemId          int     `gorm:"column:agreement_item_id;size:30;not null;primaryKey" json:"agreement_item_id"`
	AgreementId              int     `gorm:"column:agreement_id;size:30;not null;" json:"agreement_id"`
	LineTypeId               int     `gorm:"column:line_type_id;size:30;not null;" json:"line_type_id"`
	AgreementItemOperationId int     `gorm:"column:agreement_item_operation_id;size:30;not null;" json:"agreement_item_operation_id"`
	DiscountPercent          float32 `gorm:"column:discount_percent;size:30;null;" json:"discount_percent"`
	MinValue                 int     `gorm:"column:min_value;size:30;null;" json:"min_value"`
	AgreementRemark          string  `gorm:"column:agreement_remark;size:50;null" json:"agreement_remark"`
}

func (*AgreementItemDetail) TableName() string {
	return CreateAgreementItemTable
}
