package transactionsparepartentities

import "time"

const TableNameSalesOrder = "trx_sales_order"

type SalesOrder struct {
	CompanyID                        int       `gorm:"column:company_id" json:"company_id"`
	SalesOrderSystemNumber           int       `gorm:"column:sales_order_system_number;primaryKey;size:30" json:"sales_order_system_number"`
	SalesOrderDocumentNumber         string    `gorm:"column:sales_order_document_number;type:nvarchar" json:"sales_order_document_number"`
	SalesOrderStatusID               int       `gorm:"column:sales_order_status_id" json:"sales_order_status_id"`
	SalesOrderDate                   time.Time `gorm:"column:sales_order_date" json:"sales_order_date"`
	SalesOrderCloseDate              time.Time `gorm:"column:sales_order_close_date" json:"sales_order_close_date"`
	SalesEstimationDocumentNumber    string    `gorm:"column:sales_estimation_document_number;type:nvarchar" json:"sales_estimation_document_number"`
	BrandID                          int       `gorm:"column:brand_id" json:"brand_id"`
	CostProfitCenterID               int       `gorm:"column:cost_profit_center_id" json:"cost_profit_center_id"`
	EventNumberID                    int       `gorm:"column:event_number_id" json:"event_number_id"`
	SalesOrderIsAffiliated           bool      `gorm:"column:sales_order_is_affiliated" json:"sales_order_is_affiliated"`
	TransactionTypeID                int       `gorm:"column:transaction_type_id" json:"transaction_type_id"`
	SalesOrderIsOneTimeCustomer      bool      `gorm:"column:sales_order_is_one_time_customer" json:"sales_order_is_one_time_customer"`
	CustomerID                       int       `gorm:"column:customer_id" json:"customer_id"`
	TermOfPaymentID                  int       `gorm:"column:term_of_payment_id" json:"term_of_payment_id"`
	SameTaxArea                      bool      `gorm:"column:same_tax_area" json:"same_tax_area"`
	EstimatedTimeOfDelivery          time.Time `gorm:"column:estimated_time_of_delivery" json:"estimated_time_of_delivery"`
	DeliveryAddressSameAsInvoice     bool      `gorm:"column:delivery_address_same_as_invoice" json:"delivery_address_same_as_invoice"`
	DeliveryContactPerson            string    `gorm:"column:delivery_contact_person;type:nvarchar" json:"delivery_contact_person"`
	DeliveryAddress                  string    `gorm:"column:delivery_address;type:nvarchar" json:"delivery_address"`
	DeliveryAddressLine1             string    `gorm:"column:delivery_address_line1;type:nvarchar" json:"delivery_address_line1"`
	DeliveryAddressLine2             string    `gorm:"column:delivery_address_line2;type:nvarchar" json:"delivery_address_line2"`
	PurchaseOrderSystemNumber        int       `gorm:"column:purchase_order_system_number" json:"purchase_order_system_number"`
	OrderTypeID                      int       `gorm:"column:order_type_id" json:"order_type_id"`
	DeliveryViaID                    int       `gorm:"column:delivery_via_id" json:"delivery_via_id"`
	PayType                          string    `gorm:"column:pay_type;type:nvarchar" json:"pay_type"`
	WarehouseGroupID                 int       `gorm:"column:warehouse_group_id" json:"warehouse_group_id"`
	BackOrder                        bool      `gorm:"column:back_order" json:"back_order"`
	SetOrder                         bool      `gorm:"column:set_order" json:"set_order"`
	DownPaymentAmount                float32   `gorm:"column:down_payment_amount" json:"down_payment_amount"`
	DownPaymentPaidAmount            float32   `gorm:"column:down_payment_paid_amount" json:"down_payment_paid_amount"`
	DownPaymentPaymentAllocated      float32   `gorm:"column:down_payment_payment_allocated" json:"down_payment_payment_allocated"`
	DownPaymentPaymentVAT            float32   `gorm:"column:down_payment_payment_vat" json:"down_payment_payment_vat"`
	DownPaymentAllocatedToInvoice    float32   `gorm:"column:down_payment_allocated_to_invoice" json:"down_payment_allocated_to_invoice"`
	DownPaymentVATAllocatedToInvoice float32   `gorm:"column:down_payment_vat_allocated_to_invoice" json:"down_payment_vat_allocated_to_invoice"`
	Remark                           string    `gorm:"column:remark;type:nvarchar" json:"remark"`
	AgreementID                      int       `gorm:"column:agreement_id" json:"agreement_id"`
	SalesEmployeeID                  int       `gorm:"column:sales_employee_id" json:"sales_employee_id"`
	CurrencyID                       int       `gorm:"column:currency_id" json:"currency_id"`
	CurrencyExchangeID               float32   `gorm:"column:currency_exchange_id" json:"currency_exchange_id"`
	CurrencyRateTypeID               int       `gorm:"column:currency_rate_type_id" json:"currency_rate_type_id"`
	Total                            float32   `gorm:"column:total" json:"total"`
	TotalDiscount                    float32   `gorm:"column:total_discount" json:"total_discount"`
	AdditionalDiscountPercentage     float32   `gorm:"column:additional_discount_percentage" json:"additional_discount_percentage"`
	AdditionalDiscountAmount         float32   `gorm:"column:additional_discount_amount" json:"additional_discount_amount"`
	AdditionalDiscountStatusID       string    `gorm:"column:additional_discount_status_id;type:char" json:"additional_discount_status_id"`
	VATTaxID                         int       `gorm:"column:vat_tax_id" json:"vat_tax_id"`
	VATTaxPercentage                 float32   `gorm:"column:vat_tax_percentage" json:"vat_tax_percentage"`
	TotalVAT                         float32   `gorm:"column:total_vat" json:"total_vat"`
	TotalAfterVAT                    float32   `gorm:"column:total_after_vat" json:"total_after_vat"`
	ApprovalRequestNumber            int       `gorm:"column:approval_request_number" json:"approval_request_number"`
	ApprovalRemark                   string    `gorm:"column:approval_remark;type:nvarchar" json:"approval_remark"`
	VehicleSalesOrderSystemNumber    int       `gorm:"column:vehicle_sales_order_system_number" json:"vehicle_sales_order_system_number"`
	VehicleSalesOrderDetailID        int       `gorm:"column:vehicle_sales_order_detail_id" json:"vehicle_sales_order_detail_id"`
	DeliveryPhoneNumber              string    `gorm:"column:delivery_phone_number;type:nvarchar" json:"delivery_phone_number"`
	PurchaseOrderTypeID              int       `gorm:"column:purchase_order_type_id" json:"purchase_order_type_id"`
	JournalSystemNumber              int       `gorm:"column:journal_system_number" json:"journal_system_number"`
	CostCenterID                     int       `gorm:"column:cost_center_id" json:"cost_center_id"`
	DownPaymentOutstandingAmount     float32   `gorm:"column:down_payment_outstanding_amount" json:"down_payment_outstanding_amount"`
	ETDTime                          string    `gorm:"column:etd_time;type:nvarchar" json:"etd_time"`
}

func (*SalesOrder) TableName() string {
	return TableNameSalesOrder
}
