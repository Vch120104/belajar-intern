package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderDetails struct {
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

type WorkOrderResponse struct {
	StatusCode int                   `json:"status_code"`
	Message    string                `json:"message"`
	Data       WorkOrderResponseData `json:"data"`
}

type WorkOrderResponseData struct {
	WorkOrderSystemNumber                int                        `json:"work_order_system_number"`
	WorkOrderDocumentNumber              string                     `json:"work_order_document_number"`
	WorkOrderDate                        string                     `json:"work_order_date"`
	WorkOrderTypeId                      int                        `json:"work_order_type_id"`
	WorkOrderTypeName                    string                     `json:"work_order_type_name"`
	WorkOrderStatusId                    int                        `json:"work_order_status_id"`
	WorkOrderStatusName                  string                     `json:"work_order_status_name"`
	ServiceAdvisorId                     int                        `json:"service_advisor_id"`
	BrandId                              int                        `json:"brand_id"`
	BrandName                            string                     `json:"brand_name"`
	ModelId                              int                        `json:"model_id"`
	ModelName                            string                     `json:"model_name"`
	VariantId                            int                        `json:"variant_id"`
	VariantName                          string                     `json:"variant_name"`
	ServiceSite                          string                     `json:"servicesite"`
	VehicleId                            int                        `json:"vehicle_id"`
	VehicleCode                          string                     `json:"vehicle_code"`
	VehicleTnkb                          string                     `json:"vehicle_tnkb"`
	CustomerId                           int                        `json:"customer_id"`
	BillToCustomerId                     int                        `json:"billto_customer_id"`
	CampaignId                           int                        `json:"campaign_id"`
	AgreementId                          int                        `json:"agreement_id"`
	BookingSystemNumber                  int                        `json:"booking_system_number"`
	EstimationSystemNumber               int                        `json:"estimation_system_number"`
	ContractSystemNumber                 int                        `json:"contract_system_number"`
	CompanyId                            int                        `json:"company_id"`
	DealerRepresentativeId               int                        `json:"dealer_representative_id"`
	FromEra                              bool                       `json:"from_era"`
	QueueSystemNumber                    int                        `json:"queue_system_number"`
	WorkOrderArrivalTime                 string                     `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage              int                        `json:"work_order_current_mileage"`
	Storing                              bool                       `json:"storing"`
	WorkOrderRemark                      string                     `json:"work_order_remark"`
	Unregistered                         bool                       `json:"unregistered"`
	ProfitCenterId                       int                        `json:"profit_center_id"`
	WorkOrderEraNo                       string                     `json:"work_order_era_no"`
	WorkOrderEraExpiredDate              string                     `json:"work_order_era_expired_date"`
	TitlePrefix                          string                     `json:"title_prefix"`
	NameCustomer                         string                     `json:"name_customer"`
	PhoneCustomer                        string                     `json:"phone_customer"`
	MobileCustomer                       string                     `json:"mobile_customer"`
	MobileCustomerAlternative            string                     `json:"mobile_customer_alternative"`
	MobileCustomerDriver                 string                     `json:"mobile_customer_driver"`
	ContactVia                           string                     `json:"contact_via"`
	InsuranceCheck                       bool                       `json:"insurance_check"`
	InsurancePolicyNo                    string                     `json:"insurance_policy_no"`
	InsuranceExpiredDate                 string                     `json:"insurance_expired_date"`
	InsuranceClaimNo                     string                     `json:"insurance_claim_no"`
	InsurancePic                         string                     `json:"insurance_pic"`
	InsuranceWorkorderNumber             string                     `json:"insurance_workorder_number"`
	InsuranceOwnRisk                     int                        `json:"insurance_own_risk"`
	EstimationDuration                   int                        `json:"estimation_duration"`
	CustomerExpress                      bool                       `json:"customer_express"`
	LeaveCar                             bool                       `json:"leave_car"`
	CarWash                              bool                       `json:"car_wash"`
	PromiseDate                          string                     `json:"promise_date"`
	PromiseTime                          string                     `json:"promise_time"`
	FsCouponNo                           string                     `json:"fs_coupon_no"`
	Notes                                string                     `json:"notes"`
	Suggestion                           string                     `json:"suggestion"`
	AdditionalDiscountStatusApproval     int                        `json:"additional_discount_status_approval"`
	AdditionalDiscountStatusApprovalDesc string                     `json:"additional_discount_status_approval_description"`
	InvoiceSystemNumber                  int                        `json:"invoice_system_number"`
	CurrencyId                           int                        `json:"currency_id"`
	CurrencyCode                         string                     `json:"currency_code"`
	AtpmWarrantyClaimFormDocumentNumber  string                     `json:"atpm_warranty_claim_form_document_number"`
	AtpmWarrantyClaimFormDate            *string                    `json:"atpm_warranty_claim_form_date"`
	AtpmFreeServiceDocumentNumber        string                     `json:"atpm_free_service_document_number"`
	AtpmFreeServiceDate                  *string                    `json:"atpm_free_service_date"`
	TotalAfterDiscount                   *float64                   `json:"total_after_discount"`
	ApprovalRequestNumber                int                        `json:"approval_request_number"`
	JournalSystemNumber                  int                        `json:"journal_system_number"`
	ApprovalGatepassRequestNumber        int                        `json:"approval_gatepass_request_number"`
	DownPaymentAmount                    float64                    `json:"downpayment_amount"`
	DownPaymentPayment                   *float64                   `json:"downpayment_payment"`
	DownPaymentAllocated                 *float64                   `json:"downpayment_payment_allocated"`
	DownPaymentVat                       *float64                   `json:"downpayment_payment_vat"`
	DownPaymentToInvoice                 *float64                   `json:"downpayment_payment_to_invoice"`
	DownPaymentVatToInvoice              *float64                   `json:"downpayment_payment_vat_to_invoice"`
	JournalOverpaySystemNumber           int                        `json:"journal_overpay_system_number"`
	DownPaymentOverpay                   *float64                   `json:"downpayment_overpay"`
	WorkOrderSiteTypeId                  int                        `json:"work_order_site_type_id"`
	CostCenterId                         int                        `json:"cost_center_id"`
	JobOnHoldReason                      string                     `json:"job_on_hold_reason"`
	ContactPersonTitlePrefix             string                     `json:"contact_person_title_prefix"`
	WorkOrderCampaign                    NestedResponse             `json:"work_order_campaign"`
	WorkOrderGeneralRepairAgreement      NestedResponse             `json:"work_order_general_repair_agreement"`
	WorkOrderBooking                     NestedResponse             `json:"work_order_booking"`
	WorkOrderEstimation                  NestedResponse             `json:"work_order_estimation"`
	WorkOrderContract                    NestedResponse             `json:"work_order_contract"`
	WorkOrderCurrentUserDetail           WorkOrderCurrentUserDetail `json:"work_order_current_user_detail"`
	WorkOrderVehicleDetail               NestedResponse             `json:"work_order_vehicle_detail"`
	WorkOrderStnkDetail                  NestedResponse             `json:"work_order_stnk_detail"`
	WorkOrderBillingDetail               NestedResponse             `json:"work_order_billing_detail"`
	WorkOrderDetailsService              NestedResponse             `json:"work_order_details_service"`
	WorkOrderDetailsVehicle              NestedResponse             `json:"work_order_details_vehicle"`
	WorkOrderDetails                     PaginatedWorkOrderDetails  `json:"work_order_details"`
}

type NestedResponse struct {
	Data []interface{} `json:"data"`
}

type PaginatedWorkOrderDetails struct {
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
	TotalRows  int                `json:"total_rows"`
	Data       []WorkOrderDetails `json:"data"`
}

type WorkOrderCurrentUserDetail struct {
	Data []WorkOrderCurrentUserData `json:"data"`
}

type WorkOrderCurrentUserData struct {
	CustomerId      int    `json:"customer_id"`
	CustomerName    string `json:"customer_name"`
	CustomerCode    string `json:"customer_code"`
	AddressId       int    `json:"address_id"`
	AddressStreet1  string `json:"address_street_1"`
	AddressStreet2  string `json:"address_street_2"`
	AddressStreet3  string `json:"address_street_3"`
	VillageId       int    `json:"village_id"`
	VillageName     string `json:"village_name"`
	DistrictId      int    `json:"district_id"`
	DistrictName    string `json:"district_name"`
	CityId          int    `json:"city_id"`
	CityName        string `json:"city_name"`
	ProvinceId      int    `json:"province_id"`
	ProvinceName    string `json:"province_name"`
	ZipCode         string `json:"zip_code"`
	CurrentUserNpwp string `json:"current_user_npwp"`
}

type UpdateWorkOrderRequest struct {
	WorkOrderStatusId int `json:"work_order_status_id"`
}

func GetWorkOrderById(id int) (WorkOrderResponse, *exceptions.BaseErrorResponse) {
	var workOrderReq WorkOrderResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &workOrderReq)
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
