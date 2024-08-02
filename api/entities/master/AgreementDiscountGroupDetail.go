package masterentities

var CreateAgreementDiscountGroupDetailTable = "mtr_agreement_discount_group_detail"

type AgreementDiscountGroupDetail struct {
	AgreementDiscountGroupId  int     `gorm:"column:agreement_discount_group_id;size:30;not null;primaryKey" json:"agreement_discount_group_id"`
	AgreementId               int     `gorm:"column:agreement_id;size:30;not null;" json:"agreement_id"`
	AgreementSelection        int     `gorm:"column:agreement_selection_id;size:30;not null" json:"agreement_selection"`
	AgreementOrderType        int     `gorm:"column:agreement_order_type_id;size:30;not null" json:"agreement_order_type"`
	AgreementDiscountMarkupId int     `gorm:"column:agreement_discount_markup_id;size:30;not null" json:"agreement_discount_markup_id"`
	AgreementDiscount         float32 `gorm:"column:agreement_discount;size:30;not null" json:"agreement_discount"`
	AgreementDetailRemarks    string  `gorm:"column:agreement_detail_remarks;size:50;not null" json:"agreement_detail_remarks"`
}

func (*AgreementDiscountGroupDetail) TableName() string {
	return CreateAgreementDiscountGroupDetailTable
}
