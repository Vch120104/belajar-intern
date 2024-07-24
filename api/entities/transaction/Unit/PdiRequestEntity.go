package transactionunitentities

import "time"

var CreatePDIRequestTable = "trx_pdi_request"

type PdiRequest struct {
	CompanyId                string    `gorm:"column:company_id;not null" json:"company_id"`
	PdiRequestSystemNumber   int       `gorm:"column:pdi_request_system_number; primarykey;size:30" json:"pdi_request_system_number"`
	PdiRequestDocumentNumber string    `gorm:"column:pdi_request_document_number;size:25" json:"pdi_request_document_number"`
	BrandId                  int       `gorm:"column:brand_id;size:30" json:"brand_id"`
	PdiRequestDate           time.Time `gorm:"column:pdi_request_date" json:"pdi_request_date"`
	IssuedById               int       `gorm:"column:issued_by_id;size:30" json:"issued_by_id"`
	ServiceDealerId          int       `gorm:"column:service_dealer_id;size:30" json:"service_dealer_id"`
	PdiRequestRemark         string    `gorm:"column:pdi_request_remark;size:256" json:"pdi_request_remark"`
	PdiRequestStatusId       int       `gorm:"pdi_request_status_id;size:30" json:"pdi_request_status_id"`
	TotalFrt                 float64   `gorm:"column:total_frt" json:"total_frt"`
}

func (*PdiRequest) TableName() string{
	return CreatePDIRequestTable
}
