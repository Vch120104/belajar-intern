package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type WorkOrderResponse struct {
	StatusCode int                   `json:"status_code"`
	Message    string                `json:"message"`
	Data       WorkOrderResponseData `json:"data"`
}

type WorkOrderResponseData struct {
	WorkOrderDetails                     WorkOrderDetails `json:"work_order_details"`
	WorkOrderSystemNumber                int              `json:"work_order_system_number"`
	WorkOrderDocumentNumber              string           `json:"work_order_document_number"`
	WorkOrderDate                        time.Time        `json:"work_order_date"`
	WorkOrderTypeId                      int              `json:"work_order_type_id"`
	WorkOrderTypeName                    string           `json:"work_order_type_name"`
	WorkOrderStatusId                    int              `json:"work_order_status_id"`
	WorkOrderStatusName                  string           `json:"work_order_status_name"`
	ServiceAdvisorId                     int              `json:"service_advisor_id"`
	BrandId                              int              `json:"brand_id"`
	BrandName                            string           `json:"brand_name"`
	ModelId                              int              `json:"model_id"`
	ModelName                            string           `json:"model_name"`
	VariantId                            int              `json:"variant_id"`
	VariantDescription                   string           `json:"variant_description"`
	ServiceSite                          string           `json:"servicesite"`
	VehicleId                            int              `json:"vehicle_id"`
	VehicleCode                          string           `json:"vehicle_code"`
	VehicleTnkb                          string           `json:"vehicle_tnkb"`
	CustomerId                           int              `json:"customer_id"`
	BillToCustomerId                     int              `json:"billto_customer_id"`
	CampaignId                           int              `json:"campaign_id"`
	AgreementId                          int              `json:"agreement_id"`
	BookingSystemNumber                  int              `json:"booking_system_number"`
	EstimationSystemNumber               int              `json:"estimation_system_number"`
	ContractSystemNumber                 int              `json:"contract_system_number"`
	CompanyId                            int              `json:"company_id"`
	DealerRepresentativeId               int              `json:"dealer_representative_id"`
	FromEra                              bool             `json:"from_era"`
	QueueSystemNumber                    int              `json:"queue_system_number"`
	WorkOrderArrivalTime                 string           `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage              int              `json:"work_order_current_mileage"`
	Storing                              bool             `json:"storing"`
	WorkOrderRemark                      string           `json:"work_order_remark"`
	Unregistered                         bool             `json:"unregistered"`
	ProfitCenterId                       int              `json:"profit_center_id"`
	WorkOrderEraNo                       string           `json:"work_order_era_no"`
	WorkOrderEraExpiredDate              string           `json:"work_order_era_expired_date"`
	EstimationDuration                   int              `json:"estimation_duration"`
	CustomerExpress                      bool             `json:"customer_express"`
	LeaveCar                             bool             `json:"leave_car"`
	CarWash                              bool             `json:"car_wash"`
	PromiseDate                          string           `json:"promise_date"`
	PromiseTime                          string           `json:"promise_time"`
	FsCouponNo                           string           `json:"fs_coupon_no"`
	Notes                                string           `json:"notes"`
	Suggestion                           string           `json:"suggestion"`
	AdditionalDiscountStatusApproval     int              `json:"additional_discount_status_approval"`
	AdditionalDiscountStatusApprovalDesc string           `json:"additional_discount_status_approval_description"`
	InvoiceSystemNumber                  int              `json:"invoice_system_number"`
	CurrencyId                           int              `json:"currency_id"`
	CurrencyCode                         string           `json:"currency_code"`
	AtpmWarrantyClaimFormDocumentNumber  string           `json:"atpm_warranty_claim_form_document_number"`
	AtpmWarrantyClaimFormDate            *string          `json:"atpm_warranty_claim_form_date"`
	AtpmFreeServiceDocumentNumber        string           `json:"atpm_free_service_document_number"`
	AtpmFreeServiceDate                  *string          `json:"atpm_free_service_date"`
	TotalAfterDiscount                   *float64         `json:"total_after_discount"`
	ApprovalRequestNumber                int              `json:"approval_request_number"`
	JournalSystemNumber                  int              `json:"journal_system_number"`
	ApprovalGatepassRequestNumber        int              `json:"approval_gatepass_request_number"`
	DownPaymentAmount                    float64          `json:"downpayment_amount"`
	DownPaymentPayment                   *float64         `json:"downpayment_payment"`
	DownPaymentAllocated                 *float64         `json:"downpayment_payment_allocated"`
	DownPaymentVat                       *float64         `json:"downpayment_payment_vat"`
	DownPaymentToInvoice                 *float64         `json:"downpayment_payment_to_invoice"`
	DownPaymentVatToInvoice              *float64         `json:"downpayment_payment_vat_to_invoice"`
	JournalOverpaySystemNumber           int              `json:"journal_overpay_system_number"`
	DownPaymentOverpay                   *float64         `json:"downpayment_overpay"`
	WorkOrderSiteTypeId                  int              `json:"work_order_site_type_id"`
	CostCenterId                         int              `json:"cost_center_id"`
	JobOnHoldReason                      string           `json:"job_on_hold_reason"`
	ContactPersonTitlePrefix             string           `json:"contact_person_title_prefix"`
	TitlePrefix                          string           `json:"title_prefix"`
	NameCustomer                         string           `json:"name_customer"`
	PhoneCustomer                        string           `json:"phone_customer"`
	MobileCustomer                       string           `json:"mobile_customer"`
	MobileCustomerAlternative            string           `json:"mobile_customer_alternative"`
	MobileCustomerDriver                 string           `json:"mobile_customer_driver"`
	ContactVia                           string           `json:"contact_via"`
	InsuranceCheck                       bool             `json:"insurance_check"`
	InsurancePolicyNo                    string           `json:"insurance_policy_no"`
	InsuranceExpiredDate                 string           `json:"insurance_expired_date"`
	InsuranceClaimNo                     string           `json:"insurance_claim_no"`
	InsurancePic                         string           `json:"insurance_pic"`
	InsuranceWorkorderNumber             string           `json:"insurance_workorder_number"`
	InsuranceOwnRisk                     int              `json:"insurance_own_risk"`
}

type WorkOrderDetailsInfo struct {
	InvoiceSystemNumber                 int     `json:"invoice_system_number"`
	WorkOrderDetailId                   int     `json:"work_order_detail_id"`
	WorkOrderSystemNumber               int     `json:"work_order_system_number"`
	LineTypeId                          int     `json:"line_type_id"`
	LineTypeCode                        string  `json:"line_type_code"`
	TransactionTypeId                   int     `json:"transaction_type_id"`
	TransactionTypeCode                 string  `json:"transaction_type_code"`
	JobTypeId                           int     `json:"job_type_id"`
	JobTypeCode                         string  `json:"job_type_code"`
	WarehouseGroupId                    int     `json:"warehouse_group_id"`
	OperationItemId                     int     `json:"operation_item_id"`
	FrtQuantity                         float64 `json:"frt_quantity"`
	SupplyQuantity                      float64 `json:"supply_quantity"`
	OperationItemPrice                  float64 `json:"operation_item_price"`
	OperationItemDiscountAmount         float64 `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float64 `json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float64 `json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float64 `json:"operation_item_discount_request_percent"`
	WorkOrderStatusId                   int     `json:"work_order_status_id"`
	TotalCostOfGoodsSoldAmount          float64 `json:"total_cost_of_goods_sold"`
	AtpmClaimNumber                     string  `json:"atpm_claim_number"`
}

type WorkOrderDetails struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
	TotalRows  int               `json:"total_rows"`
	Data       []WorkOrderDetail `json:"data"`
}

type WorkOrderDetail struct {
	WorkOrderDetailId                   int       `json:"work_order_detail_id"`
	WorkOrderSystemNumber               int       `json:"work_order_system_number"`
	WorkOrderStatusId                   int       `json:"work_order_status_id"`
	WorkOrderStatusName                 string    `json:"work_order_status_name"`
	LineTypeId                          int       `json:"line_type_id"`
	LineTypeCode                        string    `json:"line_type_code"`
	LineTypeName                        string    `json:"line_type_name"`
	TransactionTypeId                   int       `json:"transaction_type_id"`
	TransactionTypeCode                 string    `json:"transaction_type_code"`
	JobTypeId                           int       `json:"job_type_id"`
	JobTypeCode                         string    `json:"job_type_code"`
	WarehouseGroupId                    int       `json:"warehouse_group_id"`
	OperationItemId                     int       `json:"operation_item_id"`
	FrtQuantity                         float64   `json:"frt_quantity"`
	SupplyQuantity                      float64   `json:"supply_quantity"`
	Description                         string    `json:"description"`
	OperationItemPrice                  float64   `json:"operation_item_price"`
	OperationItemDiscountAmount         float64   `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float64   `json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float64   `json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float64   `json:"operation_item_discount_request_percent"`
	OperationItemCode                   string    `json:"operation_item_code"`
	WarrantyClaimTypeId                 int       `json:"warranty_claim_type_id"`
	WarrantyClaimTypeDescription        string    `json:"warranty_claim_type_description"`
	TotalCostOfGoodsSold                float64   `json:"total_cost_of_goods_sold"`
	ServiceCategoryId                   int       `json:"service_category_id"`
	PphAmount                           float64   `json:"pph_amount"`
	TaxId                               int       `json:"tax_id"`
	PphTaxRate                          float64   `json:"pph_tax_rate"`
	LastApprovalBy                      string    `json:"last_approval_by"`
	LastApprovalDate                    time.Time `json:"last_approval_date"`
	QualityControlStatus                string    `json:"quality_control_status"`
	QualityControlExtraFrt              float64   `json:"quality_control_extra_frt"`
	QualityControlExtraReason           string    `json:"quality_control_extra_reason"`
	SubstituteTypeId                    int       `json:"substitute_type_id"`
	SubstituteTypeDescription           string    `json:"substitute_type_description"`
	AtpmClaimNumber                     string    `json:"atpm_claim_number"`
	AtpmClaimDate                       time.Time `json:"atpm_claim_date"`
	FreeServiceDocumentNumber           string    `json:"free_service_document_number"`
	FreeServiceDate                     time.Time `json:"free_service_date"`
}

type UpdateWorkOrderRequest struct {
	WorkOrderStatusId int `json:"work_order_status_id"`
}

func GetWorkOrderById(id int) (WorkOrderResponseData, *exceptions.BaseErrorResponse) {
	var workOrderReq WorkOrderResponseData
	url := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(id)
	fmt.Println(url)
	err := utils.CallAPI("GET", url, nil, &workOrderReq)
	fmt.Println(workOrderReq)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "work order service is temporarily unavailable"
		}

		return workOrderReq, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting work order by ID"),
		}
	}
	return workOrderReq, nil
}
