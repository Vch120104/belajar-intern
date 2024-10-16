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
type BinningListInsertPayloads struct {
	BinningDocumentStatusId     int        `json:"binning_document_status_id"`
	CompanyId                   int        `json:"company_id"`
	BinningDocumentNumber       string     `json:"binning_document_number"`
	BinningDocumentDate         *time.Time `json:"binning_document_date"`
	BinningReferenceTypeId      int        `json:"binning_reference_type_id"`
	ReferenceSystemNumber       int        `json:"reference_system_number"`
	ReferenceDocumentNumber     string     `json:"reference_document_number"`
	WarehouseGroupId            int        `json:"warehouse_group_id"`
	WarehouseId                 int        `json:"warehouse_id"`
	SupplierId                  int        `json:"supplier_id"`
	SupplierDeliveryOrderNumber string     `json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string     `json:"supplier_invoice_number"`
	SupplierInvoiceDate         *time.Time `json:"supplier_invoice_date"`
	SupplierFakturPajakNumber   string     `json:"supplier_faktur_pajak_number"`
	SupplierFakturPajakDate     *time.Time `json:"supplier_faktur_pajak_date"`
	SupplierDeliveryPerson      string     `json:"supplier_delivery_person"`
	SupplierCaseNumber          string     `json:"supplier_case_number"`
	CurrencyId                  int        `json:"currency_id "`
	ExchangeId                  int        `json:"exchange_id "`
	CreatedByUserId             int        `json:"created_by_user_id"`
	CreatedDate                 *time.Time `json:"created_date"`
	UpdatedByUserId             int        `json:"updated_by_user_id"`
	UpdatedDate                 *time.Time `json:"updated_date"`
	ChangeNo                    int        `json:"change_no"`
	BinningTypeId               int        `json:"binning_type_id"`
}
type BinningListSavePayload struct {
	BinningSystemNumber         int        `json:"binning_system_number"`
	BinningDocumentStatusId     int        `json:"binning_document_status_id"`
	CompanyId                   int        `json:"company_id"`
	BinningDocumentNumber       string     `json:"binning_document_number"`
	BinningDocumentDate         *time.Time `json:"binning_document_date"`
	BinningReferenceTypeId      int        `json:"binning_reference_type_id"`
	ReferenceSystemNumber       int        `json:"reference_system_number"`
	ReferenceDocumentNumber     string     `json:"reference_document_number"`
	WarehouseGroupId            int        `json:"warehouse_group_id"`
	WarehouseId                 int        `json:"warehouse_id"`
	SupplierId                  int        `json:"supplier_id"`
	SupplierDeliveryOrderNumber string     `json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string     `json:"supplier_invoice_number"`
	SupplierInvoiceDate         *time.Time `json:"supplier_invoice_date"`
	SupplierFakturPajakNumber   string     `json:"supplier_faktur_pajak_number"`
	SupplierFakturPajakDate     *time.Time `json:"supplier_faktur_pajak_date"`
	SupplierDeliveryPerson      string     `json:"supplier_delivery_person"`
	SupplierCaseNumber          string     `json:"supplier_case_number"`
	CurrencyId                  int        `json:"currency_id "`
	ExchangeId                  int        `json:"exchange_id "`
	UpdatedByUserId             int        `json:"updated_by_user_id"`
	UpdatedDate                 *time.Time `json:"updated_date"`
	ChangeNo                    int        `json:"change_no"`
	BinningTypeId               int        `json:"binning_type_id"`
}
type BinningListGetByIdResponses struct {
	BinningSystemNumber         int     `gorm:"column:binning_system_number;not null;size:30" json:"binning_system_number"`
	BinningLineNumber           int     `gorm:"column:binning_line_number;not null" json:"binning_line_number"`
	ItemCode                    string  `gorm:"column:item_code;not null;size:30" json:"item_code"`
	ItemName                    string  `json:"item_name"`
	UomCode                     string  `json:"uom_code"`
	WarehouseLocationCode       string  `json:"warehouse_location_code"`
	ItemPrice                   float64 `gorm:"column:item_price;not null" json:"item_price"`
	PurchaseOrderQuantity       int     `gorm:"column:purchase_order_quantity;not null" json:"purchase_order_quantity"`
	DeliveryOrderQuantity       int     `gorm:"column:delivery_order_quantity;not null" json:"delivery_order_quantity"`
	ReferenceSystemNumber       int     `gorm:"column:reference_system_number;not null" json:"reference_system_number"`
	ReferenceLineNumber         int     `gorm:"column:reference_line_number;not null" json:"reference_line_number"`
	OriginalItemCode            string  `gorm:"column:original_item_code;null" json:"original_item_code"`
	PurchaseOrderDocumentNumber string  `gorm:"column:purchase_order_document_number" json:"purchase_order_document_number"`
}

type BinningListDetailInsertPayloads struct {
	BinningSystemNumber         int        `gorm:"column:binning_system_number;not null;size:30" json:"binning_system_number"`
	BinningLineNumber           int        `gorm:"column:binning_line_number;not null" json:"binning_line_number"`
	ItemId                      int        `gorm:"column:item_id;not null;size:30" json:"item_id"`
	UomId                       int        `gorm:"column:uom_id;not null;size:30" json:"uom_id"`
	ItemPrice                   float64    `gorm:"column:item_price;not null" json:"item_price"`
	WarehouseLocationId         int        `gorm:"column:warehouse_location_id;not null;size:30" json:"warehouse_location_id"`
	PurchaseOrderQuantity       int        `gorm:"column:purchase_order_quantity;not null" json:"purchase_order_quantity"`
	DeliveryOrderQuantity       float64    `gorm:"column:delivery_order_quantity;not null" json:"delivery_order_quantity"`
	ReferenceSystemNumber       int        `gorm:"column:reference_system_number;not null" json:"reference_system_number"`
	ReferenceLineNumber         int        `gorm:"column:reference_line_number;not null" json:"reference_line_number"`
	ReferenceDetailSystemNumber int        `gorm:"column:reference_detail_system_number;not null" json:"reference_detail_system_number"`
	GoodsReceiveSystemNumber    int        `gorm:"column:goods_receive_system_number;not null" json:"goods_receive_system_number"`
	GoodsReceiveLineNumber      int        `gorm:"column:goods_receive_line_number;null" json:"goods_receive_line_number"`
	OriginalItemId              int        `gorm:"column:original_item_id;null" json:"original_item_id"`
	CreatedByUserId             int        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                 *time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId             int        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                 *time.Time `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo                    int        `gorm:"column:change_no;size:30;" json:"change_no"`
}
