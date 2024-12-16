package transactionsparepartentities

import "time"

const TableNameSalesOrder = "trx_sales_order"

type SalesOrder struct {
	CompanyID                        int        `gorm:"column:company_id" json:"company_id"`
	SalesOrderSystemNumber           int        `gorm:"column:sales_order_system_number;primaryKey;size:30" json:"sales_order_system_number"`
	SalesOrderDocumentNumber         string     `gorm:"column:sales_order_document_number;size:25" json:"sales_order_document_number"`
	SalesOrderStatusID               int        `gorm:"column:sales_order_status_id" json:"sales_order_status_id"`
	SalesOrderDate                   time.Time  `gorm:"column:sales_order_date" json:"sales_order_date"`
	SalesOrderCloseDate              time.Time  `gorm:"column:sales_order_close_date" json:"sales_order_close_date"`
	SalesEstimationDocumentNumber    string     `gorm:"column:sales_estimation_document_number;size:25" json:"sales_estimation_document_number"`
	BrandID                          int        `gorm:"column:brand_id" json:"brand_id"`
	ProfitCenterID                   int        `gorm:"column:profit_center_id;size:30" json:"profit_center_id"`
	EventNumberID                    int        `gorm:"column:event_number_id;size:30" json:"event_number_id"`
	SalesOrderIsAffiliated           bool       `gorm:"column:sales_order_is_affiliated" json:"sales_order_is_affiliated"`
	TransactionTypeID                int        `gorm:"column:transaction_type_id;size:30" json:"transaction_type_id"`
	SalesOrderIsOneTimeCustomer      bool       `gorm:"column:sales_order_is_one_time_customer" json:"sales_order_is_one_time_customer"`
	CustomerID                       int        `gorm:"column:customer_id;size:30" json:"customer_id"`
	TermOfPaymentID                  int        `gorm:"column:term_of_payment_id;size:30" json:"term_of_payment_id"`
	SameTaxArea                      bool       `gorm:"column:same_tax_area" json:"same_tax_area"`
	EstimatedTimeOfDelivery          time.Time  `gorm:"column:estimated_time_of_delivery" json:"estimated_time_of_delivery"`
	DeliveryAddressSameAsInvoice     bool       `gorm:"column:delivery_address_same_as_invoice" json:"delivery_address_same_as_invoice"`
	DeliveryContactPerson            string     `gorm:"column:delivery_contact_person;size:40" json:"delivery_contact_person"`
	DeliveryAddress                  string     `gorm:"column:delivery_address;size:100" json:"delivery_address"`
	DeliveryAddressLine1             string     `gorm:"column:delivery_address_line1;size:100" json:"delivery_address_line1"`
	DeliveryAddressLine2             string     `gorm:"column:delivery_address_line2;size:100" json:"delivery_address_line2"`
	PurchaseOrderSystemNumber        int        `gorm:"column:purchase_order_system_number" json:"purchase_order_system_number"`
	OrderTypeID                      int        `gorm:"column:order_type_id" json:"order_type_id"`
	DeliveryViaID                    int        `gorm:"column:delivery_via_id" json:"delivery_via_id"`
	PayType                          string     `gorm:"column:pay_type;size:15" json:"pay_type"`
	WarehouseGroupID                 int        `gorm:"column:warehouse_group_id" json:"warehouse_group_id"`
	BackOrder                        bool       `gorm:"column:back_order" json:"back_order"`
	SetOrder                         bool       `gorm:"column:set_order" json:"set_order"`
	DownPaymentAmount                float64    `gorm:"column:down_payment_amount" json:"down_payment_amount"`
	DownPaymentPaidAmount            float64    `gorm:"column:down_payment_paid_amount" json:"down_payment_paid_amount"`
	DownPaymentPaymentAllocated      float64    `gorm:"column:down_payment_payment_allocated" json:"down_payment_payment_allocated"`
	DownPaymentPaymentVAT            float64    `gorm:"column:down_payment_payment_vat" json:"down_payment_payment_vat"`
	DownPaymentAllocatedToInvoice    float64    `gorm:"column:down_payment_allocated_to_invoice" json:"down_payment_allocated_to_invoice"`
	DownPaymentVATAllocatedToInvoice float64    `gorm:"column:down_payment_vat_allocated_to_invoice" json:"down_payment_vat_allocated_to_invoice"`
	Remark                           string     `gorm:"column:remark;size:100" json:"remark"`
	AgreementID                      int        `gorm:"column:agreement_id" json:"agreement_id"`
	SalesEmployeeID                  int        `gorm:"column:sales_employee_id" json:"sales_employee_id"`
	CurrencyID                       int        `gorm:"column:currency_id" json:"currency_id"`
	CurrencyExchangeID               float64    `gorm:"column:currency_exchange_id" json:"currency_exchange_id"`
	CurrencyRateTypeID               int        `gorm:"column:currency_rate_type_id" json:"currency_rate_type_id"`
	Total                            float64    `gorm:"column:total" json:"total"`
	TotalDiscount                    float64    `gorm:"column:total_discount" json:"total_discount"`
	AdditionalDiscountPercentage     float64    `gorm:"column:additional_discount_percentage" json:"additional_discount_percentage"`
	AdditionalDiscountAmount         float64    `gorm:"column:additional_discount_amount" json:"additional_discount_amount"`
	AdditionalDiscountStatusID       int        `gorm:"column:additional_discount_status_id;size:30" json:"additional_discount_status_id"`
	VATTaxID                         int        `gorm:"column:vat_tax_id" json:"vat_tax_id"`
	VATTaxPercentage                 float64    `gorm:"column:vat_tax_percentage" json:"vat_tax_percentage"`
	TotalVAT                         float64    `gorm:"column:total_vat" json:"total_vat"`
	TotalAfterVAT                    float64    `gorm:"column:total_after_vat" json:"total_after_vat"`
	ApprovalRequestNumber            int        `gorm:"column:approval_request_number" json:"approval_request_number"`
	ApprovalRemark                   string     `gorm:"column:approval_remark;size:50" json:"approval_remark"`
	VehicleSalesOrderSystemNumber    int        `gorm:"column:vehicle_sales_order_system_number" json:"vehicle_sales_order_system_number"`
	VehicleSalesOrderDetailID        int        `gorm:"column:vehicle_sales_order_detail_id" json:"vehicle_sales_order_detail_id"`
	DeliveryPhoneNumber              string     `gorm:"column:delivery_phone_number;size:25" json:"delivery_phone_number"`
	PurchaseOrderTypeID              int        `gorm:"column:purchase_order_type_id" json:"purchase_order_type_id"`
	JournalSystemNumber              int        `gorm:"column:journal_system_number" json:"journal_system_number"`
	CostCenterID                     int        `gorm:"column:cost_center_id" json:"cost_center_id"`
	DownPaymentOutstandingAmount     float64    `gorm:"column:down_payment_outstanding_amount" json:"down_payment_outstanding_amount"`
	ETDTime                          time.Time  `gorm:"column:etd_time;nvarchar" json:"etd_time"`
	AtpmInternalPurpose              string     `gorm:"atpm_internal_purpose;size:15" json:"atpm_internal_purpose"`
	CreatedByUserId                  int        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                      *time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId                  int        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                      *time.Time `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo                         int        `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*SalesOrder) TableName() string {
	return TableNameSalesOrder
}
