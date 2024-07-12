package transactionsparepartpayloads

import "time"

type SalesOrderResponse struct {
	CompanyID                         int       `json:"company_id"`
	SalesOrderSystemNumber            int       `json:"sales_order_system_number"`
	SalesOrderDocumentNumber          string    `json:"sales_order_document_number"`
	SalesOrderStatusID                int       `json:"sales_order_status_id"`
	SalesOrderDate                    time.Time `json:"sales_order_date"`
	SalesOrderCloseDate               time.Time `json:"sales_order_close_date"`
	SalesEstimationDocumentNumber     string    `json:"sales_estimation_document_number"`
	BrandID                           int       `json:"brand_id"`
	CostProfitCenterID                int       `json:"cost_profit_center_id"`
	EventNumberID                     int       `json:"event_number_id"`
	SalesOrderIsAffiliated            bool      `json:"sales_order_is_affiliated"`
	TransactionTypeID                 int       `json:"transaction_type_id"`
	SalesOrderIsOneTimeCustomer       bool      `json:"sales_order_is_one_time_customer"`
	CustomerID                        int       `json:"customer_id"`
	TermOfPaymentID                   int       `json:"term_of_payment_id"`
	SameTaxArea                       bool      `json:"same_tax_area"`
	EstimatedTimeOfDelivery           time.Time `json:"estimated_time_of_delivery"`
	DeliveryAddressSameAsInvoice      bool      `json:"delivery_address_same_as_invoice"`
	DeliveryContactPerson             string    `json:"delivery_contact_person"`
	DeliveryAddress                   string    `json:"delivery_address"`
	DeliveryAddressLine1              string    `json:"delivery_address_line1"`
	DeliveryAddressLine2              string    `json:"delivery_address_line2"`
	PurchaseOrderSystemNumber         int       `json:"purchase_order_system_number"`
	OrderTypeID                       int       `json:"order_type_id"`
	DeliveryViaID                     int       `json:"delivery_via_id"`
	PayType                           string    `json:"pay_type"`
	WarehouseGroupID                  int       `json:"warehouse_group_id"`
	BackOrder                         bool      `json:"back_order"`
	SetOrder                          bool      `json:"set_order"`
	DownPaymentAmount                 float32   `json:"down_payment_amount"`
	DownPaymentPaidAmount             float32   `json:"down_payment_paid_amount"`
	DownPaymentPaymentAllocated       float32   `json:"down_payment_payment_allocated"`
	DownPaymentPaymentVAT             float32   `json:"down_payment_payment_vat"`
	DownPaymentAllocatedToInvoice     float32   `json:"down_payment_allocated_to_invoice"`
	DownPaymentVATAllocatedToInvoice  float32   `json:"down_payment_vat_allocated_to_invoice"`
	Remark                            string    `json:"remark"`
	AgreementID                       int       `json:"agreement_id"`
	SalesEmployeeID                   int       `json:"sales_employee_id"`
	CurrencyID                        int       `json:"currency_id"`
	CurrencyExchangeID                float32   `json:"currency_exchange_id"`
	CurrencyRateTypeID                int       `json:"currency_rate_type_id"`
	Total                             float32   `json:"total"`
	TotalDiscount                     float32   `json:"total_discount"`
	TotalDiscountAmount               float32   `json:"total_discount_amount"`
	TotalDiscountAmountVAT            float32   `json:"total_discount_amount_vat"`
	TotalDiscountAmountAllocated      float32   `json:"total_discount_amount_allocated"`
	TotalDiscountAmountVATAllocated   float32   `json:"total_discount_amount_vat_allocated"`
	TotalDiscountAmountAllocatedTo    float32   `json:"total_discount_amount_allocated_to"`
	TotalDiscountAmountVATAllocatedTo float32   `json:"total_discount_amount_vat_allocated_to"`
	TotalDiscountAmountAllocatedToVAT float32   `json:"total_discount_amount_allocated_to_vat"`
}
