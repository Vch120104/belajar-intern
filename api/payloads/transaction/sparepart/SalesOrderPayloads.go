package transactionsparepartpayloads

import (
	"time"
)

type SupplyTypeResponse struct {
	SupplyTypeId          int    `json:"supply_type_id"`
	SupplyTypeDescription string `json:"supply_type_description"`
}

type SalesOrderResponse struct {
	CompanyID                         int       `json:"company_id"`
	SalesOrderSystemNumber            int       `json:"sales_order_system_number"`
	SalesOrderDocumentNumber          string    `json:"sales_order_document_number"`
	SalesOrderStatusID                int       `json:"sales_order_status_id"`
	SalesOrderDate                    time.Time `json:"sales_order_date"`
	SalesOrderCloseDate               time.Time `json:"sales_order_close_date"`
	SalesEstimationDocumentNumber     string    `json:"sales_estimation_document_number"`
	BrandID                           int       `json:"brand_id"`
	ProfitCenterID                    int       `json:"profit_center_id"`
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
	DownPaymentAmount                 float64   `json:"down_payment_amount"`
	DownPaymentPaidAmount             float64   `json:"down_payment_paid_amount"`
	DownPaymentPaymentAllocated       float64   `json:"down_payment_payment_allocated"`
	DownPaymentPaymentVAT             float64   `json:"down_payment_payment_vat"`
	DownPaymentAllocatedToInvoice     float64   `json:"down_payment_allocated_to_invoice"`
	DownPaymentVATAllocatedToInvoice  float64   `json:"down_payment_vat_allocated_to_invoice"`
	Remark                            string    `json:"remark"`
	AgreementID                       int       `json:"agreement_id"`
	SalesEmployeeID                   int       `json:"sales_employee_id"`
	CurrencyID                        int       `json:"currency_id"`
	CurrencyExchangeID                float64   `json:"currency_exchange_id"`
	CurrencyRateTypeID                int       `json:"currency_rate_type_id"`
	Total                             float64   `json:"total"`
	TotalDiscount                     float64   `json:"total_discount"`
	TotalDiscountAmount               float64   `json:"total_discount_amount"`
	TotalDiscountAmountVAT            float64   `json:"total_discount_amount_vat"`
	TotalDiscountAmountAllocated      float64   `json:"total_discount_amount_allocated"`
	TotalDiscountAmountVATAllocated   float64   `json:"total_discount_amount_vat_allocated"`
	TotalDiscountAmountAllocatedTo    float64   `json:"total_discount_amount_allocated_to"`
	TotalDiscountAmountVATAllocatedTo float64   `json:"total_discount_amount_vat_allocated_to"`
	TotalDiscountAmountAllocatedToVAT float64   `json:"total_discount_amount_allocated_to_vat"`
}

type SalesOrderInsertHeaderPayload struct {
	CompanyID                     int       `json:"company_id"`
	SalesOrderDate                time.Time `json:"sales_order_date"`
	SalesEstimationDocumentNumber string    `json:"sales_estimation_document_number"`
	BrandID                       int       `json:"brand_id"`
	ProfitCenterID                int       `json:"profit_center_id"`
	EventNumberID                 int       `json:"event_number_id"`
	SalesOrderIsAffiliated        bool      `json:"sales_order_is_affiliated"`
	TransactionTypeID             int       `json:"transaction_type_id"`
	SalesOrderIsOneTimeCustomer   bool      `json:"sales_order_is_one_time_customer"`
	CustomerID                    int       `json:"customer_id"`
	TermOfPaymentID               int       `json:"term_of_payment_id"`
	SameTaxArea                   bool      `json:"same_tax_area"`
	EstimatedTimeOfDelivery       time.Time `json:"estimated_time_of_delivery"`
	DeliveryAddressSameAsInvoice  bool      `json:"delivery_address_same_as_invoice"`
	DeliveryContactPerson         string    `json:"delivery_contact_person"`
	DeliveryAddress               string    `json:"delivery_address"`
	DeliveryAddressLine1          string    `json:"delivery_address_line1"`
	DeliveryAddressLine2          string    `json:"delivery_address_line2"`
	DeliveryPhoneNumber           string    `json:"delivery_phone_number"`
	PurchaseOrderSystemNumber     int       `json:"purchase_order_system_number"`
	PurchaseOrderTypeID           int       `json:"purchase_order_type_id"`
	DeliveryViaID                 int       `json:"delivery_via_id"`
	PayType                       string    `json:"pay_type"`
	WarehouseGroupID              int       `json:"warehouse_group_id"`
	BackOrder                     bool      `json:"back_order"`
	SetOrder                      bool      `json:"set_order"`
	DownPaymentAmount             float64   `json:"down_payment_amount"`
	Remark                        string    `json:"remark"`
	AgreementID                   int       `json:"agreement_id"`
	SalesEmployeeID               int       `json:"sales_employee_id"`
	CurrencyID                    int       `json:"currency_id"`
	CurrencyExchangeID            float64   `json:"currency_exchange_id"`
	CurrencyRateTypeID            int       `json:"currency_rate_type_id"`
	Total                         float64   `json:"total"`
	TotalDiscount                 float64   `json:"total_discount"`
	AdditionalDiscountPercentage  float64   `json:"additional_discount_percentage"`
	AdditionalDiscountAmount      float64   `json:"additional_discount_amount"`
	AdditionalDiscountStatusID    int       `json:"additional_discount_status_id"`
	VATTaxPercentage              float64   `json:"vat_tax_percentage"`
	VATTaxId                      int       `json:"vat_tax_id"`
	TotalVAT                      float64   `json:"total_vat"`
	TotalAfterVAT                 float64   `json:"total_after_vat"`
	VehicleSalesOrderSystemNumber int       `json:"vehicle_sales_order_system_number"`
	VehicleSalesOrderDetailID     int       `json:"vehicle_sales_order_detail_id"`
	OrderTypeID                   int       `json:"order_type_id"`
	CostCenterID                  int       `json:"cost_center_id"`
	ETDTime                       time.Time `json:"etd_time"`
	CreatedByUserId               int       `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
}

type SalesOrderEstimationGetByIdResponse struct {
	CompanyID                           int       `gorm:"column:company_id" json:"company_id"`
	SalesOrderSystemNumber              int       `gorm:"column:sales_order_system_number;primaryKey;size:30" json:"sales_order_system_number"`
	SalesOrderDocumentNumber            string    `gorm:"column:sales_order_document_number;type:nvarchar" json:"sales_order_document_number"`
	SalesOrderStatusID                  int       `gorm:"column:sales_order_status_id" json:"sales_order_status_id"`
	SalesOrderStatusDescription         string    `json:"sales_order_status_description"`
	SalesOrderDate                      time.Time `gorm:"column:sales_order_date" json:"sales_order_date"`
	SalesEstimationDocumentNumber       string    `gorm:"column:sales_estimation_document_number;type:nvarchar" json:"sales_estimation_document_number"`
	BrandID                             int       `gorm:"column:brand_id" json:"brand_id"`
	VehicleBrandCode                    string    `json:"vehicle_brand_code"`
	VehicleBrandDescription             string    `json:"vehicle_brand_description"`
	ProfitCenterID                      int       `gorm:"column:profit_center_id;size:30" json:"profit_center_id"`
	EventNumberID                       int       `gorm:"column:event_number_id;size:30" json:"event_number_id"`
	SalesOrderIsAffiliated              bool      `gorm:"column:sales_order_is_affiliated" json:"sales_order_is_affiliated"`
	TransactionTypeID                   int       `gorm:"column:transaction_type_id;size:30" json:"transaction_type_id"`
	TransactionTypeDescription          string    `json:"transaction_type_description"`
	SalesOrderIsOneTimeCustomer         bool      `gorm:"column:sales_order_is_one_time_customer" json:"sales_order_is_one_time_customer"`
	CustomerID                          int       `gorm:"column:customer_id;size:30" json:"customer_id"`
	TermOfPaymentID                     int       `gorm:"column:term_of_payment_id;size:30" json:"term_of_payment_id"`
	TermOfPaymentDescription            string    `json:"term_of_payment_description"`
	SameTaxArea                         bool      `gorm:"column:same_tax_area" json:"same_tax_area"`
	EstimatedTimeOfDelivery             time.Time `gorm:"column:estimated_time_of_delivery" json:"estimated_time_of_delivery"`
	DeliveryAddressSameAsInvoice        bool      `gorm:"column:delivery_address_same_as_invoice" json:"delivery_address_same_as_invoice"`
	DeliveryContactPerson               string    `gorm:"column:delivery_contact_person;type:nvarchar" json:"delivery_contact_person"`
	DeliveryAddress                     string    `gorm:"column:delivery_address;type:nvarchar" json:"delivery_address"`
	DeliveryAddressLine1                string    `gorm:"column:delivery_address_line1;type:nvarchar" json:"delivery_address_line1"`
	DeliveryAddressLine2                string    `gorm:"column:delivery_address_line2;type:nvarchar" json:"delivery_address_line2"`
	PurchaseOrderSystemNumber           int       `gorm:"column:purchase_order_system_number" json:"purchase_order_system_number"`
	PurchaseOrdeCompanyId               int       `json:"purchase_orde_company_id"`
	PurchaseOrderCompanyCode            string    `json:"purchase_order_company_code"`
	PurchaseOrderDocumentNumber         string    `json:"purchase_order_document_number"`
	OrderTypeID                         int       `gorm:"column:order_type_id" json:"order_type_id"`
	DeliveryViaID                       int       `gorm:"column:delivery_via_id" json:"delivery_via_id"`
	DeliveryViaDescription              string    `json:"delivery_via_description"`
	PayType                             string    `gorm:"column:pay_type;type:nvarchar" json:"pay_type"`
	WarehouseGroupID                    int       `gorm:"column:warehouse_group_id" json:"warehouse_group_id"`
	WarehouseGroupDescription           string    `gorm:"column:warehouse_group_description" json:"warehouse_group_description"`
	BackOrder                           bool      `gorm:"column:back_order" json:"back_order"`
	SetOrder                            bool      `gorm:"column:set_order" json:"set_order"`
	DownPaymentAmount                   float64   `gorm:"column:down_payment_amount" json:"down_payment_amount"`
	DownPaymentPaidAmount               float64   `gorm:"column:down_payment_paid_amount" json:"down_payment_paid_amount"`
	DownPaymentPaymentAllocated         float64   `gorm:"column:down_payment_payment_allocated" json:"down_payment_payment_allocated"`
	DownPaymentPaymentVAT               float64   `gorm:"column:down_payment_payment_vat" json:"down_payment_payment_vat"`
	DownPaymentAllocatedToInvoice       float64   `gorm:"column:down_payment_allocated_to_invoice" json:"down_payment_allocated_to_invoice"`
	DownPaymentVATAllocatedToInvoice    float64   `gorm:"column:down_payment_vat_allocated_to_invoice" json:"down_payment_vat_allocated_to_invoice"`
	Remark                              string    `gorm:"column:remark;type:nvarchar" json:"remark"`
	AgreementID                         int       `gorm:"column:agreement_id" json:"agreement_id"`
	SalesEmployeeID                     int       `gorm:"column:sales_employee_id" json:"sales_employee_id"`
	CurrencyID                          int       `gorm:"column:currency_id" json:"currency_id"`
	CurrencyExchangeID                  float64   `gorm:"column:currency_exchange_id" json:"currency_exchange_id"`
	CurrencyRateTypeID                  int       `gorm:"column:currency_rate_type_id" json:"currency_rate_type_id"`
	Total                               float64   `gorm:"column:total" json:"total"`
	TotalDiscount                       float64   `gorm:"column:total_discount" json:"total_discount"`
	AdditionalDiscountPercentage        float64   `gorm:"column:additional_discount_percentage" json:"additional_discount_percentage"`
	AdditionalDiscountAmount            float64   `gorm:"column:additional_discount_amount" json:"additional_discount_amount"`
	AdditionalDiscountStatusID          int       `gorm:"column:additional_discount_status_id;size:30" json:"additional_discount_status_id"`
	AdditionalDiscountStatusDescription string    `json:"additional_discount_status_description"`
	VATTaxID                            int       `gorm:"column:vat_tax_id" json:"vat_tax_id"`
	VATTaxPercentage                    float64   `gorm:"column:vat_tax_percentage" json:"vat_tax_percentage"`
	TotalVAT                            float64   `gorm:"column:total_vat" json:"total_vat"`
	TotalAfterVAT                       float64   `gorm:"column:total_after_vat" json:"total_after_vat"`
	ApprovalRequestNumber               int       `gorm:"column:approval_request_number" json:"approval_request_number"`
	ApprovalRemark                      string    `gorm:"column:approval_remark;type:nvarchar" json:"approval_remark"`
	VehicleSalesOrderSystemNumber       int       `gorm:"column:vehicle_sales_order_system_number" json:"vehicle_sales_order_system_number"`
	VehicleSalesOrderDetailID           int       `gorm:"column:vehicle_sales_order_detail_id" json:"vehicle_sales_order_detail_id"`
	PurchaseOrderTypeID                 int       `gorm:"column:purchase_order_type_id" json:"purchase_order_type_id"`
	CostCenterID                        int       `gorm:"column:cost_center_id" json:"cost_center_id"`
	IsAtpm                              bool      `json:"is_atpm"`
	AtpmInternalPurpose                 string    `json:"atpm_internal_purpose"`
	//AtpmInternalPurpose
}

type GetAllSalesOrderResponse struct {
	SalesOrderSystemNumber        int       `gorm:"column:sales_order_system_number;primaryKey;size:30" json:"sales_order_system_number"`
	SalesOrderDocumentNumber      string    `gorm:"column:sales_order_document_number;type:nvarchar" json:"sales_order_document_number"`
	SalesEstimationDocumentNumber string    `gorm:"column:sales_estimation_document_number;type:nvarchar" json:"sales_estimation_document_number"`
	SalesOrderDate                time.Time `json:"sales_order_date"`
	ReferenceDocumentNumber       string    `json:"reference_document_number"`
	CustomerID                    int       `gorm:"column:customer_id;size:30" json:"customer_id"`
	CustomerName                  string    `json:"customer_name"`
	SalesOrderStatusID            int       `json:"sales_order_status_id"`
	SalesOrderStatusDescription   string    `json:"sales_order_status_description"`
	TransactionTypeID             int       `gorm:"column:transaction_type_id;size:30" json:"transaction_type_id"`
	TransactionTypeDescription    string    `json:"transaction_type_description"`
	CreatedByUserId               int       `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedByUserName             string    `json:"created_by_user_name"`
	PurchaseOrderSystemNumber     int       `gorm:"column:purchase_order_system_number" json:"purchase_order_system_number"`
	VehicleSalesOrderSystemNumber int       `gorm:"column:vehicle_sales_order_system_number" json:"vehicle_sales_order_system_number"`
}

// validate:"required"`
type SalesOrderDetailInsertPayload struct {
	SalesOrderSystemNumber              int        `json:"sales_order_system_number" validate:"required"`
	SalesOrderLineStatusId              *int       `json:"sales_order_line_status_id"` // FK to mtr_approval_status in general-service
	SalesOrderDetailSystemNumber        int        `json:"sales_order_detail_system_number"`
	ItemSubstituteId                    *int       `json:"item_substitute_id"`
	ItemId                              int        `json:"item_id"`
	QuantityDemand                      float64    `json:"quantity_demand" validate:"required"`
	IsAvailable                         bool       `json:"is_available" validate:"required"`
	QuantitySupply                      float64    `json:"quantity_supply" validate:"required"`
	QuantityPick                        float64    `json:"quantity_pick" validate:"required"`
	UomId                               *int       `json:"uom_id"`
	Price                               float64    `json:"price" validate:"required"`
	PriceEffectiveDate                  *time.Time `json:"price_effective_date"`
	DiscountPercent                     float64    `json:"discount_percent"`
	DiscountRequestPercent              float64    `json:"discount_request_percent"`
	Remark                              string     `json:"remark"`
	ApprovalRequestNumber               *int       `json:"approval_request_number"` // FK to trx_approval_request_source in ?
	ApprovalRemark                      string     `json:"approval_remark"`
	VehicleSalesOrderSystemNumber       *int       `json:"vehicle_sales_order_system_number"`        // FK to trx_vehicle_sales_order in sales-service
	VehicleSalesOrderDetailSystemNumber *int       `json:"vehicle_sales_order_detail_system_number"` // FK to trx_vehicle_sales_order_detail in sales-service
	PriceListId                         *int       `json:"price_list_id"`
	AdditionalDiscountPercentage        float64    `json:"additional_discount_percentage"`
	HeaderRemark                        string     `json:"header_remark"`
}

type SalesOrderDeleteDetailResponse struct {
	DeleteMessage string `json:"delete_message"`
	DeleteStatus  bool   `json:"delete_status"`
}

//
//type SalesOrderProposedDiscountResponse struct {
//	ProposedDiscount float64 `json:"proposed_discount" validate:"required"`
//}
