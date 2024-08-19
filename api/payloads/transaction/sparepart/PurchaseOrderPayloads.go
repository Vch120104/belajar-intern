package transactionsparepartpayloads

import "time"

type GetAllDBResponses struct {
	PurchaseOrderSystemNumber int `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order"`
	//WarehouseId int `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatusId       int        `json:"purchase_order_status_id" parent_entity:"trx_item_purchase_order"`
	OrderTypeId                 int        `json:"order_type_id" parent_entity:"trx_item_purchase_order"`
	WarehouseId                 int        `json:"warehouse_id" parent_entity:"trx_item_purchase_order"`
	SupplierId                  int        `json:"supplier_id" parent_entity:"trx_item_purchase_order"`
	PurchaseRequestSystemNumber int        `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
}

type GetAllResponses struct {
	PurchaseOrderSystemNumber int `json:"purchase_order_system_number" parent_entity:"trx_item_purchase_order"`
	//WarehouseId int `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	PurchaseOrderDocumentNumber string     `json:"purchase_order_document_number" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderDocumentDate   *time.Time `json:"purchase_order_document_date" parent_entity:"trx_item_purchase_order"`
	PurchaseOrderStatus         string     `json:"purchase_order_status" parent_entity:"trx_item_purchase_order"`
	OrderType                   string     `json:"order_type" parent_entity:"trx_item_purchase_order"`
	WarehouseName               string     `json:"warehouse_name" parent_entity:"trx_item_purchase_order"`
	SupplierName                string     `json:"supplier_name" parent_entity:"trx_item_purchase_order"`
	PurchaseRequestSystemNumber string     `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
}

//
//type PurchaseOrderEntities struct {
//	CompanyId                           int        `gorm:"column:company_id;size:30;" json:"company_id"`
//	PurchaseOrderSystemNumber           int        `gorm:"column:purchase_order_system_number;size:30;not null;primaryKey;" json:"purchase_order_system_number"`
//	PurchaseOrderDocumentNumber         string     `gorm:"column:purchase_order_document_number;size:30;" json:"purchase_order_document_number"`
//	PurchaseOrderDocumentDate           *time.Time `gorm:"column:purchase_order_document_date;size:30;" json:"purchase_order_document_date"`
//	PurchaseOrderStatusId               int        `gorm:"column:purchase_order_status_id;size:30;" json:"purchase_order_status_id"`
//	BrandId                             int        `gorm:"column:brand_id;size:30;" json:"brand_id"`
//	ItemGroupId                         int        `gorm:"column:item_group_id;size:30;" json:"item_group_id"`
//	OrderTypeId                         int        `gorm:"column:order_type_id;size:30;" json:"order_type_id"`
//	SupplierId                          int        `gorm:"column:supplier_id;size:30;" json:"supplier_id"`
//	SupplierPicId                       int        `gorm:"column:supplier_pic_id;size:30;" json:"supplier_pic_id"`
//	WarehouseId                         int        `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
//	CostCenterId                        int        `gorm:"column:cost_center_id;size:2;" json:"cost_center_id"`
//	ProfitType                          string     `gorm:"column:profit_type;size:30;" json:"profit_type"`
//	ProfitCenterId                      int        `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
//	AffiliatedPurchaseOrder             bool       `gorm:"column:affiliated_purchase_order" json:"affiliated_purchase_order"`
//	CurrencyId                          int        `gorm:"column:currency_id;size:30;" json:"currency_id"`
//	BackOrder                           bool       `gorm:"column:back_order;" json:"back_order"`
//	SetOrder                            bool       `gorm:"set_order;" json:"set_order"`
//	ViaBinning                          bool       `gorm:"via_binning;" json:"via_binning"`
//	VatCode                             string     `gorm:"column:vat_code;size:30;" json:"vat_code"`
//	TotalDiscount                       *float64   `gorm:"column:total_discount;" json:"total_discount"`
//	TotalAmount                         *float64   `gorm:"column:total_amount;" json:"total_amount"`
//	TotalVat                            *float64   `gorm:"column:total_vat;" json:"total_vat"`
//	TotalAfterVat                       *float64   `gorm:"column:total_after_vat;" json:"total_after_vat"`
//	LastTotalDiscount                   *float64   `gorm:"column:last_total_discount;" json:"last_total_discount"`
//	LastTotalAmount                     *float64   `gorm:"column:last_total_amount;" json:"last_total_amount"`
//	LastTotalVat                        *float64   `gorm:"column:last_total_vat;" json:"last_total_vat"`
//	LastTotalAfterVat                   *float64   `gorm:"column:last_total_after_vat;" json:"last_total_after_vat"`
//	TotalAmountConfirm                  *float64   `gorm:"column:total_amount_confirm;" json:"total_amount_confirm"`
//	PurchaseOrderRemark                 string     `gorm:"column:purchase_order_remark;size:256;" json:"purchase_order_remark"`
//	DpRequest                           *float64   `gorm:"column:dp_request;" json:"dp_request"`
//	DpPayment                           *float64   `gorm:"column:dp_payment;" json:"dp_payment"`
//	DpPaymentAllocated                  *float64   `gorm:"column:dp_payment_allocated;" json:"dp_payment_allocated"`
//	DpPaymentAllocatedInvoice           *float64   `gorm:"column:dp_payment_allocated_invoice;" json:"dp_payment_allocated_invoice"`
//	DpPaymentAllocatedPpn               *float64   `gorm:"column:dp_payment_allocated_ppn;" json:"dp_payment_allocated_ppn"`
//	DpPaymentAllocatedRequestForPayment *float64   `gorm:"column:dp_payment_allocated_request_for_payment;" json:"dp_payment_allocated_request_for_payment"`
//	DeliveryId                          int        `gorm:"column:delivery_id;" json:"delivery_id"`
//	ExpectedDeliveryDate                *time.Time `gorm:"column:expected_delivery_date;" json:"expected_delivery_date"`
//	ExpectedArrivalDate                 *time.Time `gorm:"column:expected_arrival_date;" json:"expected_arrival_date"`
//	EstimatedDeliveryDate               *time.Time `gorm:"column:estimated_delivery_date;" json:"estimated_delivery_date"`
//	EstimatedDeliveryTime               string     `gorm:"column:estimated_delivery_time;size:5;" json:"estimated_delivery_time"`
//	SalesOrderSystemNumber              int        `gorm:"column:sales_order_system_number;" json:"sales_order_system_number"`
//	SalesOrderDocumentNumber            string     `gorm:"column:sales_order_document_number;size:25;" json:"sales_order_document_number"`
//	LastPrintById                       int        `gorm:"column:last_print_by_id;" json:"last_print_by_id"`
//	ApprovalRequestById                 int        `gorm:"column:approval_request_by_id;" json:"approval_request_by_id"`
//	ApprovalRequestNumber               int        `gorm:"column:approval_request_number;" json:"approval_request_number"`
//	ApprovalRequestDate                 *time.Time `gorm:"column:approval_request_date;" json:"approval_request_date"`
//	ApprovalRemark                      string     `gorm:"column:approval_remark;size:256;" json:"approval_remark"`
//	ApprovalLastById                    int        `gorm:"column:approval_last_by_id;" json:"approval_last_by_id"`
//	ApprovalLastDate                    *time.Time `gorm:"column:approval_last_date;" json:"approval_last_date"`
//	TotalInvoiceDownPayment             *float64   `gorm:"column:total_invoice_down_payment;" json:"total_invoice_down_payment"`
//	TotalInvoiceDownPaymentVat          *float64   `gorm:"column:total_invoice_down_payment_vat;" json:"total_invoice_down_payment_vat"`
//	TotalInvoiceDownPaymentAfterVat     *float64   `gorm:"column:total_invoice_down_payment_after_vat;" json:"total_invoice_down_payment_after_vat"`
//	DownPaymentReturn                   *float64   `gorm:"column:down_payment_return;" json:"down_payment_return"`
//	JournalSystemNumber                 int        `gorm:"column:journal_system_number;" json:"journal_system_number"`
//	EventNumber                         string     `gorm:"column:event_number;size:10;" json:"event_number"`
//	ItemClassId                         int        `gorm:"column:item_class_id;" json:"item_class_id"`
//	IsDirectShipment                    string     `gorm:"column:is_direct_shipment;size:1;" json:"is_direct_shipment"`
//	CustomerId                          int        `gorm:"column:customer_id;" json:"customer_id"`
//	ExternalPurchaseOrderNumber         string     `gorm:"column:external_purchase_order_number;size:10;" json:"external_purchase_order_number"`
//	PurchaseOrderTypeId                 int        `gorm:"column:purchase_order_type_id;" json:"purchase_order_type_id"`
//	CurrencyExchangeRate                *float64   `gorm:"column:currency_exchange_rate;" json:"currency_exchange_rate"`
//	//PurchaseOrderDetail                 []PurchaseOrderDetailEntities `gorm:"foreignKey:PurchaseOrderSystemNumber;references:PurchaseOrderSystemNumber" json:"work_order_detail"`
//}
