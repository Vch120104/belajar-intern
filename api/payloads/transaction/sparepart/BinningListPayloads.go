package transactionsparepartpayloads

import "time"

type BinningListGetByIdResponse struct {
	CompanyId                   int        `json:"company_id"`
	BinningSystemNumber         int        `json:"binning_system_number"`
	BinningDocumentStatusId     int        `json:"binning_document_status_id"`
	BinningDocumentNumber       string     `json:"binning_document_number"`
	BinningDocumentDate         *time.Time `json:"binning_document_date"`
	BinningReferenceType        string     `json:"binning_reference_type_id"`
	ReferenceSystemNumber       int        `json:"reference_system_number"`
	ReferenceDocumentNumber     string     `json:"reference_document_number"`
	WarehouseGroupCode          string     `json:"warehouse_group_code"`
	WarehouseCode               string     `json:"warehouse_code"`
	SupplierCode                string     `json:"supplier_code"`
	SupplierDeliveryOrderNumber string     `json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string     `json:"supplier_invoice_number"`
	SupplierInvoiceDate         *time.Time `json:"supplier_invoice_date"`
	SupplierFakturPajakNumber   string     `json:"supplier_faktur_pajak_number"`
	SupplierFakturPajakDate     *time.Time `json:"supplier_faktur_pajak_date"`
	SupplierDeliveryPerson      string     `json:"supplier_delivery_person"`
	SupplierCaseNumber          string     `json:"supplier_case_number"`
	ItemGroup                   string     `json:"item_group"`
	CreatedByUserId             int        `json:"created_by_user_id"`
	CreatedDate                 *time.Time `json:"created_date"`
	UpdatedByUserId             int        `json:"updated_by_user_id"`
	UpdatedDate                 *time.Time `json:"updated_date"`
	ChangeNo                    int        `json:"change_no"`
}

type BinningListGetPaginationResponse struct {
	BinningSystemNumber         int        `json:"binning_system_number"`
	BinningDocumentStatusId     int        `json:"binning_document_status_id"`
	BinningDocumentNumber       string     `json:"binning_document_number"`
	BinningDocumentDate         *time.Time `json:"binning_document_date"`
	ReferenceDocumentNumber     string     `json:"reference_document_number"`
	SupplierInvoiceNumber       string     `json:"supplier_invoice_number"`
	SupplierName                string     `json:"supplier_name"`
	SupplierCaseNumber          string     `json:"supplier_case_number"`
	Status                      string     `json:"status"`
	SupplierDeliveryOrderNumber string     `json:"supplier_delivery_order_number"`
}
