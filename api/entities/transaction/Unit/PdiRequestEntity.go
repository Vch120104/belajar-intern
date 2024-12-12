package transactionunitentities

import "time"

var CreatePDIRequestTable = "dms_microservices_sales_dev.dbo.trx_pdi_request"

type PdiRequest struct {
	CompanyId                int       `gorm:"column:company_id;not null" json:"company_id"`                                               // Kolom wajib
	PdiRequestSystemNumber   int       `gorm:"column:pdi_request_system_number;primaryKey;autoIncrement" json:"pdi_request_system_number"` // Kolom wajib, auto increment
	PdiRequestDocumentNumber string    `gorm:"column:pdi_request_document_number;size:25;unique" json:"pdi_request_document_number"`       // Kolom unik
	BrandId                  int       `gorm:"column:brand_id;not null" json:"brand_id"`                                                   // Kolom wajib
	PdiRequestDate           time.Time `gorm:"column:pdi_request_date" json:"pdi_request_date"`                                            // Kolom opsional
	IssuedById               int       `gorm:"column:issued_by_id;not null" json:"issued_by_id"`                                           // Kolom wajib
	ServiceDealerId          *int      `gorm:"column:service_dealer_id" json:"service_dealer_id"`                                          // Kolom opsional (nullable)
	ServiceById              int       `gorm:"column:service_by_id;not null" json:"service_by_id"`                                         // Kolom wajib
	PdiRequestRemark         *string   `gorm:"column:pdi_request_remark;size:256" json:"pdi_request_remark"`                               // Kolom opsional (nullable)
	PdiRequestStatusId       int       `gorm:"column:pdi_request_status_id;not null" json:"pdi_request_status_id"`                         // Kolom wajib
	TotalFrt                 *float64  `gorm:"column:total_frt" json:"total_frt"`                                                          // Kolom opsional (nullable)
}

func (*PdiRequest) TableName() string {
	return CreatePDIRequestTable
}
