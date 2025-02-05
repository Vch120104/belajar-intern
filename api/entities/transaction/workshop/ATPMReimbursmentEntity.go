package transactionworkshopentities

import "time"

const TableNameAtpmReimbursement = "trx_atpm_reimbursement"

type AtpmReimbursement struct {
	ClaimSystemNumber        int       `gorm:"column:claim_system_number;size:30;primaryKey" json:"claim_system_number"`
	CompanyId                int       `gorm:"column:company_id;size:30" json:"company_id"`
	ReimbursementStatusId    int       `gorm:"column:reimbursement_status_id;size:30" json:"reimbursement_status_id"`
	InvoiceSystemNumber      int       `gorm:"column:invoice_system_number;size:30" json:"invoice_system_number"`
	InvoiceDocumentNumber    string    `gorm:"column:invoice_document_number;size:30" json:"invoice_document_number"`
	InvoiceDate              time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	TaxInvoiceSystemNumber   int       `gorm:"column:tax_invoice_system_number;size:30" json:"tax_invoice_system_number"`
	TaxInvoiceDocumentNumber string    `gorm:"column:tax_invoice_document_number;size:30" json:"tax_invoice_document_number"`
	TaxInvoiceDate           time.Time `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	KwitansiSystemNumber     int       `gorm:"column:kwitansi_system_number;size:30" json:"kwitansi_system_number"`
	KwitansiDocumentNumber   string    `gorm:"column:kwitansi_document_number;size:30" json:"kwitansi_document_number"`
}

func (*AtpmReimbursement) TableName() string {
	return TableNameAtpmReimbursement
}
