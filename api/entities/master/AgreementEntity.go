package masterentities

import "time"

var CreateAgreementTable = "mtr_agreement"

type Agreement struct {
	AgreementId       int       `gorm:"column:agreement_id;size:30;not null;primaryKey" json:"agreement_id"`
	AgreementCode     string    `gorm:"column:agreement_code;size:50;not null" json:"agreement_code"`
	IsActive          bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	BrandId           int       `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	CustomerId        int       `gorm:"column:customer_id;size:30;not null" json:"customer_id"`
	ProfitCenterId    int       `gorm:"column:profit_center_id;size:30;not null" json:"profit_center_id"`
	AgreementDateFrom time.Time `gorm:"column:agreement_date_from;not null;type:datetime" json:"agreement_date_from"`
	AgreementDateTo   time.Time `gorm:"column:agreement_date_to;not null;type:datetime" json:"agreement_date_to"`
	DealerId          int       `gorm:"column:company_id;size:30;not null" json:"company_id"`
	TopId             int       `gorm:"column:top_id;size:30;not null" json:"top_id"`
	AgreementRemark   string    `gorm:"column:agreement_remark;size:50;not null" json:"agreement_remark"`
}

func (*Agreement) TableName() string {
	return CreateAgreementTable
}
