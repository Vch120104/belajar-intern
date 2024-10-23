package transactionsparepartpayloads

import "time"

type GoodsReceivesGetAllPayloads struct {
	GoodsReceiveSystemNumber    int       `json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber  string    `json:"goods_receive_document_number"`
	GoodsReceiveDocumentDate    time.Time `json:"goods_receive_document_date"`
	ItemGroupName               string    `json:"item_group_name"`
	ReferenceDocumentNumber     string    `json:"reference_document_number"`
	SupplierId                  int       `json:"supplier_id"`
	SupplierName                string    `json:"supplier_name"`
	GoodsReceiveStatusId        int       `json:"goods_receive_status_id"`
	JournalSystemNumber         int       `json:"journal_system_number"`
	SupplierDeliveryOrderNumber string    `json:"supplier_delivery_order_number"`
	QuantityGoodsReceive        float64   `json:"quantity_goods_receive"`
	TotalAmount                 float64   `json:"total_amount"`
}

type GoodsReceivesGetByIdResponses struct {
	GoodsReceiveSystemNumber    int       `json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber  string    `json:"goods_receive_document_number"`
	ItemGroupId                 int       `json:"item_group_id"`
	GoodsReceiveDocumentDate    time.Time `json:"goods_receive_document_date"`
	SupplierId                  int       `json:"supplier_id"`
	GoodsReceiveStatusId        int       `json:"goods_receive_status_id"`
	ReferenceTypeGoodReceiveId  int       `json:"reference_type_good_receive_id"`
	ReferenceSystemNumber       int       `json:"reference_system_number"`
	ReferenceDocumentNumber     string    `json:"reference_document_number"`
	AffiliatedPurchaseOrder     bool      `json:"affiliated_purchase_order"`
	ViaBinning                  bool      `json:"via_binning"`
	SetOrder                    bool      `json:"set_order"`
	BackOrder                   bool      `json:"back_order"`
	BrandId                     int       `json:"brand_id"`
	CostCenterId                int       `json:"cost_center_id"`
	ProfitCenterId              int       `json:"profit_center_id"`
	TransactionTypeId           int       `json:"transaction_type_id"`
	EventId                     int       `json:"event_id"`
	SupplierDeliveryOrderNumber string    `json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string    `json:"supplier_invoice_number"`
	SupplierTaxInvoiceNumber    string    `json:"supplier_tax_invoice_number"`
	WarehouseId                 int       `json:"warehouse_id"`
	WarehouseCode               string    `json:"warehouse_code"`
	WarehouseName               string    `json:"warehouse_name"`
	WarehouseClaimId            int       `json:"warehouse_claim_id"`
	WarehouseClaimCode          string    `json:"warehouse_claim_code"`
	WarehouseClaimName          string    `json:"warehouse_claim_name"`
	ItemClassId                 int       `json:"item_class_id"`
}
