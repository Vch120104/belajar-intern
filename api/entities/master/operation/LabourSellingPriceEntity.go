package masteroperationentities

import "time"

var CreateLabourSellingPriceTable = "mtr_labour_selling_price"

type LabourSellingPrice struct {
	IsActive             bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	LabourSellingPriceId int       `gorm:"column:labour_selling_price_id;size:30;not null;primaryKey" json:"labour_selling_price_id"`
	CompanyId            int       `gorm:"column:company_id;size:30;not null"        json:"company_id"` // Fk with mtr_company on general sevice
	BrandId              int       `gorm:"column:brand_id;size:30;not null" json:"brand_id"`            //Fk with mtr_brand on sales service
	JobTypeId            int       `gorm:"column:job_type_id;size:30;not null" json:"job_type_id"`      // Fk with mtr_job_type on general service
	EffectiveDate        time.Time `gorm:"column:effective_date;not null" json:"effective_date"`
	BillToId             int       `gorm:"column:bill_to_id;size:30;" json:"bill_to_id"`
	Description          string    `gorm:"column:description;size:128;null;" json:"description"`
}

func (*LabourSellingPrice) TableName() string {
	return CreateLabourSellingPriceTable
}
