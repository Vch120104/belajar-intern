package transactionworkshopentities

import "time"

const TableNameContractService = "trx_contract_service"

type ContractService struct {
	CompanyId                     int       `gorm:"column:company_id;size:30;not null" json:"company_id"`
	ContractServiceSystemNumber   int       `gorm:"column:contract_service_system_number;size:30;primary key" json:"contract_service_system_number"`
	ContractSevriceDocumentNumber string    `gorm:"column:contract_Service_document_number;size:256" json:"contract_service_document_number"`
	ContractServiceDate           time.Time `gorm:"column:contract_service_date;not null" json:"contract_service_date"`
	ContractServiceName           string    `gorm:"column:contract_service_name;size:40;not null" json:"contract_service_name"`
	ContractServiceFrom           time.Time `gorm:"column:contract_service_from;not null" json:"contract_service_from"`
	ContractServiceTo             time.Time `gorm:"column:contract_service_to;not null" json:"contract_service_to"`
	ContractServiceStatusId       int       `gorm:"column:contract_service_status_id;size:30;not null" json:"contract_service_status_id"`
	ProfitCenterId                int       `gorm:"column:profit_center_id;size:30;not null" json:"profit_center_id"`
	BrandId                       int       `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	ModelId                       int       `gorm:"column:model_id;size:30;not null" json:"model_id"`
	VehicleId                     int       `gorm:"column:vehicle_id;size:30;not null" json:"vehicle_id"`
	RegisteredMileage             int       `gorm:"column:registered_mileage;size:30;not null" json:"registered_mileage"`
	Remark                        string    `gorm:"column:remark;size:50" json:"remark"`
	Total                         float64   `gorm:"column:total;not null" json:"total"`
	TotalValueAfterTax            float64   `gorm:"column:total_value_after_tax;not null" json:"total_value_after_tax"`
	TaxId                         int       `gorm:"column:tax_id;size:30;not null" json:"tax_id"`
	ValueAfterTaxrate             float64   `gorm:"column:value_after_tax_rate;not null" json:"value_after_tax_rate"`
	TotalPayment                  float64   `gorm:"column:total_payment;not null" json:"total_payment"`
	TotalPaymentAllocated         float64   `gorm:"column:total_payment_allocated;not null" json:"total_payment_allocated"`
}

func (*ContractService) TableName() string {
	return TableNameContractService
}