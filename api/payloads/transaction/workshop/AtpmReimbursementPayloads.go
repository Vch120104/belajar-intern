package transactionworkshoppayloads

import "time"

type AtpmReimbursementRequest struct {
	ClaimSystemNumber        int       `json:"claim_system_number"`
	CompanyId                int       `json:"company_id"`
	ReimbursementStatusId    int       `json:"reimbursement_status_id"`
	InvoiceSystemNumber      int       `json:"invoice_system_number"`
	InvoiceDocumentNumber    string    `json:"invoice_document_number"`
	InvoiceDate              time.Time `json:"invoice_date"`
	TaxInvoiceSystemNumber   int       `json:"tax_invoice_system_number"`
	TaxInvoiceDocumentNumber string    `json:"tax_invoice_document_number"`
	TaxInvoiceDate           time.Time `json:"tax_invoice_date"`
	KwitansiSystemNumber     int       `json:"kwitansi_system_number"`
	KwitansiDocumentNumber   string    `json:"kwitansi_document_number"`
}

type AtpmReimbursementUpdate struct {
	ClaimSystemNumber        int       `json:"claim_system_number"`
	CompanyId                int       `json:"company_id"`
	ReimbursementStatusId    int       `json:"reimbursement_status_id"`
	InvoiceSystemNumber      int       `json:"invoice_system_number"`
	InvoiceDocumentNumber    string    `json:"invoice_document_number"`
	InvoiceDate              time.Time `json:"invoice_date"`
	TaxInvoiceSystemNumber   int       `json:"tax_invoice_system_number"`
	TaxInvoiceDocumentNumber string    `json:"tax_invoice_document_number"`
	TaxInvoiceDate           time.Time `json:"tax_invoice_date"`
	KwitansiSystemNumber     int       `json:"kwitansi_system_number"`
	KwitansiDocumentNumber   string    `json:"kwitansi_document_number"`
}
