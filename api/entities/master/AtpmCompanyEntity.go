package masterentities

const TableName = "mtr_atpm_company_mapping"

type AtpmCompanyMapping struct {
	IsActive             bool `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	AtpmCompanyMappingId int  `gorm:"column:atpm_company_mapping_id;size:30;not null;primaryKey" json:"atpm_company_mapping_id"`
	CompanyId            int  `gorm:"column:company_id;size:30" json:"company_id"`
	BrandId              int  `gorm:"column:brand_id;size:30;not null;primaryKey" json:"brand_id"`
}

func (*AtpmCompanyMapping) TableName() string {
	return TableName
}
