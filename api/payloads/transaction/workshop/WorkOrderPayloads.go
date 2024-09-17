package transactionworkshoppayloads

import (
	"encoding/json"
	"time"
)

type WorkOrderResponse struct {
	WorkOrderStatusId             int       `json:"work_order_status_id"`
	WorkOrderSystemNumber         int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber       string    `json:"work_order_document_number"`
	WorkOrderDate                 time.Time `json:"work_order_date"`
	WorkOrderTypeId               int       `json:"work_order_type_id"`
	ServiceAdvisorId              int       `json:"service_advisor_id"`
	WorkOrderServiceSite          int       `json:"work_order_servicesite_id"`
	BrandId                       int       `json:"brand_id"`
	ModelId                       int       `json:"model_id"`
	VariantId                     int       `json:"variant_id"`
	VehicleId                     int       `json:"vehicle_id"`
	CustomerId                    int       `json:"customer_id"`
	BilltoCustomerId              int       `json:"billto_customer_id"`
	WorkOrderStatusEra            bool      `json:"work_order_status_era"`
	WorkOrderEraNo                string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate       time.Time `json:"work_order_era_expired_date"`
	QueueSystemNumber             int       `json:"queue_system_number"`
	WorkOrderArrivalTime          time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage       int       `json:"work_order_current_mileage"`
	WorkOrderStatusStoring        bool      `json:"work_order_status_storing"`
	WorkOrderRemark               string    `json:"work_order_remark"`
	WorkOrderStatusUnregistered   bool      `json:"work_order_status_unregistered"`
	WorkOrderProfitCenter         string    `json:"work_order_profit_center"`
	WorkOrderDealerRepCode        string    `json:"work_order_dealer_rep_code"`
	CampaignId                    int       `json:"campaign_id"`
	AgreementId                   int       `json:"agreement_id"`
	ServiceRequestSystemNumber    int       `json:"system_request_system_number"`
	EstimationSystemNumber        int       `json:"estimation_system_number"`
	ContractSystemNumber          int       `json:"contract_system_number"`
	CompanyId                     int       `json:"company_id"`
	DealerRepresentativeId        int       `json:"dealer_representative_id"`
	Titleprefix                   string    `json:"title_prefix"`
	NameCust                      string    `json:"name_customer"`
	PhoneCust                     string    `json:"phone_customer"`
	MobileCust                    string    `json:"mobile_customer"`
	MobileCustAlternative         string    `json:"mobile_customer_alternative"`
	MobileCustDriver              string    `json:"mobile_customer_driver"`
	ContactVia                    string    `json:"contact_via"`
	WorkOrderStatusInsurance      bool      `json:"work_order_status_insurance"`
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`
}

type WorkOrderRequest struct {
	// Basic information
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	ServiceAdvisorId        int       `json:"service_advisor_id"`
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VariantId               int       `json:"variant_id"`
	ServiceSite             string    `json:"servicesite"`
	VehicleId               int       `json:"vehicle_id"`
	CustomerId              int       `json:"customer_id"`
	BilltoCustomerId        int       `json:"billto_customer_id"`
	CampaignId              int       `json:"campaign_id"`
	AgreementId             int       `json:"agreement_id"`
	BoookingId              int       `json:"booking_system_number"`
	EstimationId            int       `json:"estimation_system_number"`
	ContractSystemNumber    int       `json:"contract_system_number"`
	CompanyId               int       `json:"company_id"`
	DealerRepresentativeId  int       `json:"dealer_representative_id"`
	FromEra                 bool      `json:"from_era"`
	QueueSystemNumber       int       `json:"queue_system_number"`
	WorkOrderArrivalTime    string    `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage int       `json:"work_order_current_mileage"`
	Storing                 bool      `json:"storing"`
	WorkOrderRemark         string    `json:"work_order_remark"`
	Unregistered            bool      `json:"unregistered"`
	WorkOrderProfitCenter   int       `json:"work_order_profit_center"`

	// Work order status and details
	WorkOrderEraNo           string `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  string `json:"work_order_era_expired_date"`
	WorkOrderStatusId        int    `json:"work_order_status_id"`
	WorkOrderStatusInsurance bool   `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsuranceCheck       bool    `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string  `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate string  `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string  `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string  `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string  `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64 `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64 `json:"estimation_duration"`
	CustomerExpress    bool    `json:"customer_express"`
	LeaveCar           bool    `json:"leave_car"`
	CarWash            bool    `json:"car_wash"`
	PromiseDate        string  `json:"promise_date"`
	PromiseTime        string  `json:"promise_time"`

	// Additional information
	FSCouponNo        string  `json:"fs_coupon_no"`
	Notes             string  `json:"notes"`
	Suggestion        string  `json:"suggestion"`
	DownpaymentAmount float64 `json:"dp_amount"`
}

type WorkOrderResponseDetail struct {
	// Basic information
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	WorkOrderTypeName       string    `json:"work_order_type_name"`
	WorkOrderStatusId       int       `json:"work_order_status_id"`
	WorkOrderStatusName     string    `json:"work_order_status_name"`
	ServiceAdvisorId        int       `json:"service_advisor_id"`
	BrandId                 int       `json:"brand_id"`
	BrandName               string    `json:"brand_name"`
	ModelId                 int       `json:"model_id"`
	ModelName               string    `json:"model_name"`
	VariantId               int       `json:"variant_id"`
	VariantName             string    `json:"variant_name"`
	ServiceSite             string    `json:"servicesite"`
	VehicleId               int       `json:"vehicle_id"`
	VehicleCode             string    `json:"vehicle_code"`
	VehicleTnkb             string    `json:"vehicle_tnkb"`
	CustomerId              int       `json:"customer_id"`
	BilltoCustomerId        int       `json:"billto_customer_id"`
	CampaignId              int       `json:"campaign_id"`
	AgreementId             int       `json:"agreement_id"`
	BoookingId              int       `json:"booking_system_number"`
	EstimationId            int       `json:"estimation_system_number"`
	ContractSystemNumber    int       `json:"contract_system_number"`
	CompanyId               int       `json:"company_id"`
	DealerRepresentativeId  int       `json:"dealer_representative_id"`
	FromEra                 bool      `json:"from_era"`
	QueueSystemNumber       int       `json:"queue_system_number"`
	WorkOrderArrivalTime    time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage int       `json:"work_order_current_mileage"`
	Storing                 bool      `json:"storing"`
	WorkOrderRemark         string    `json:"work_order_remark"`
	Unregistered            bool      `json:"unregistered"`
	WorkOrderProfitCenterId int       `json:"profit_center_id"`

	// Work order status and details
	WorkOrderEraNo           string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  time.Time `json:"work_order_era_expired_date"`
	WorkOrderStatusInsurance bool      `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsuranceCheck       bool      `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64   `json:"estimation_duration"`
	CustomerExpress    bool      `json:"customer_express"`
	LeaveCar           bool      `json:"leave_car"`
	CarWash            bool      `json:"car_wash"`
	PromiseDate        time.Time `json:"promise_date"`
	PromiseTime        time.Time `json:"promise_time"`

	// Additional information
	FSCouponNo             string                          `json:"fs_coupon_no"`
	Notes                  string                          `json:"notes"`
	Suggestion             string                          `json:"suggestion"`
	DownpaymentAmount      float64                         `json:"dp_amount"`
	WorkOrderDetailService WorkOrderDetailsResponseRequest `json:"work_order_details_service"`
	WorkOrderDetailVehicle WorkOrderDetailsResponseVehicle `json:"work_order_details_vehicle"`
	WorkOrderDetails       WorkOrderDetailsResponse        `json:"work_order_details"`
}

type WorkOrderDetailsResponseVehicle struct {
	DataVehicle []WorkOrderServiceVehicleResponse `json:"data"`
}

type WorkOrderDetailsResponseRequest struct {
	DataRequest []WorkOrderServiceResponse `json:"data"`
}

type WorkOrderDetailsResponse struct {
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
	TotalRows  int                       `json:"total_rows"`
	Data       []WorkOrderDetailResponse `json:"data"`
}

type WorkOrderNormalRequest struct {
	// Basic information
	WorkOrderTypeId            int       `json:"work_order_type_id"`
	BookingSystemNumber        int       `json:"booking_system_number"`
	EstimationSystemNumber     int       `json:"estimation_system_number"`
	ServiceRequestSystemNumber int       `json:"service_request_system_number"`
	PDISystemNumber            int       `json:"pdi_system_number"`
	RepeatedSystemNumber       int       `json:"repeated_system_number"`
	BrandId                    int       `json:"brand_id"`
	ModelId                    int       `json:"model_id"`
	VariantId                  int       `json:"variant_id"`
	ServiceSite                string    `json:"servicesite"`
	VehicleId                  int       `json:"vehicle_id"`
	CustomerId                 int       `json:"customer_id"`
	BilltoCustomerId           int       `json:"billto_customer_id"`
	CampaignId                 int       `json:"campaign_id"`
	CompanyId                  int       `json:"company_id"`
	FromEra                    bool      `json:"from_era"`
	QueueSystemNumber          int       `json:"queue_system_number"`
	WorkOrderArrivalTime       time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage    int       `json:"work_order_current_mileage"`
	Storing                    bool      `json:"storing"`
	WorkOrderRemark            string    `json:"work_order_remark"`
	Unregistered               bool      `json:"unregistered"`
	WorkOrderProfitCenter      int       `json:"work_order_profit_center"`
	DealerRepresentativeId     int       `json:"dealer_representative_id"`

	// Work order status and details
	WorkOrderEraNo           string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  time.Time `json:"work_order_era_expired_date"`
	WorkOrderStatusId        int       `json:"work_order_status_id"`
	WorkOrderStatusInsurance bool      `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsuranceCheck       bool      `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64   `json:"estimation_duration"`
	CustomerExpress    bool      `json:"customer_express"`
	LeaveCar           bool      `json:"leave_car"`
	CarWash            bool      `json:"car_wash"`
	PromiseDate        time.Time `json:"promise_date"`
	PromiseTime        time.Time `json:"promise_time"`

	// Additional information
	FSCouponNo        string  `json:"fs_coupon_no"`
	Notes             string  `json:"notes"`
	Suggestion        string  `json:"suggestion"`
	DownpaymentAmount float64 `json:"dp_amount"`
}

type WorkOrderNormalSaveRequest struct {
	// Basic information
	BilltoCustomerId        int       `json:"billto_customer_id"`
	CampaignId              int       `json:"campaign_id"`
	CompanyId               int       `json:"company_id"`
	FromEra                 bool      `json:"from_era"`
	QueueSystemNumber       int       `json:"queue_system_number"`
	WorkOrderArrivalTime    time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage int       `json:"work_order_current_mileage"`
	Storing                 bool      `json:"storing"`
	WorkOrderRemark         string    `json:"work_order_remark"`
	Unregistered            bool      `json:"unregistered"`
	WorkOrderProfitCenter   int       `json:"work_order_profit_center"`
	DealerRepresentativeId  int       `json:"dealer_representative_id"`

	// Work order status and details
	WorkOrderEraNo           string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  time.Time `json:"work_order_era_expired_date"`
	WorkOrderStatusId        int       `json:"work_order_status_id"`
	WorkOrderStatusInsurance bool      `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsuranceCheck       bool      `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64   `json:"estimation_duration"`
	CustomerExpress    bool      `json:"customer_express"`
	LeaveCar           bool      `json:"leave_car"`
	CarWash            bool      `json:"car_wash"`
	PromiseDate        time.Time `json:"promise_date"`
	PromiseTime        time.Time `json:"promise_time"`

	// Additional information
	FSCouponNo        string  `json:"fs_coupon_no"`
	Notes             string  `json:"notes"`
	Suggestion        string  `json:"suggestion"`
	DownpaymentAmount float64 `json:"dp_amount"`
}

type WorkOrderGetAllRequest struct {

	// Basic information
	WorkOrderSystemNumber   int       `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number" parent_entity:"trx_work_order"`
	WorkOrderDate           time.Time `json:"work_order_date" parent_entity:"trx_work_order"`
	WorkOrderTypeId         int       `json:"work_order_type_id" parent_entity:"trx_work_order"`
	ServiceAdvisorId        int       `json:"service_advisor_id" parent_entity:"trx_work_order"`
	BrandId                 int       `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId                 int       `json:"model_id" parent_entity:"trx_work_order"`
	VariantId               int       `json:"variant_id" parent_entity:"trx_work_order"`
	ServiceSite             string    `json:"service_site" parent_entity:"trx_work_order"`
	VehicleId               int       `json:"vehicle_id" parent_entity:"trx_work_order"`
	CustomerId              int       `json:"customer_id" parent_entity:"trx_work_order"`
	BilltoCustomerId        int       `json:"billable_to_id" parent_entity:"trx_work_order"`
	StatusId                int       `json:"work_order_status_id" parent_entity:"trx_work_order"`
	RepeatedJob             int       `json:"repeated_system_number" parent_entity:"trx_work_order"`
}

type WorkOrderGetAllResponse struct {
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	FormattedWorkOrderDate  string    `json:"formatted_work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	WorkOrderTypeName       string    `json:"work_order_type_name"`
	StatusId                int       `json:"work_order_status_id"`
	StatusName              string    `json:"work_order_status_description"`
	ServiceAdvisorId        int       `json:"service_advisor_id"`
	BrandId                 int       `json:"brand_id"`
	BrandName               string    `json:"brand_name"`
	ModelId                 int       `json:"model_id"`
	ModelName               string    `json:"model_name"`
	VariantId               int       `json:"variant_id"`
	ServiceSite             string    `json:"service_site"`
	VehicleId               int       `json:"vehicle_id"`
	VehicleCode             string    `json:"vehicle_chassis_number"`
	VehicleTnkb             string    `json:"vehicle_registration_certificate_tnkb"`
	CustomerId              int       `json:"customer_id"`
	BilltoCustomerId        int       `json:"billable_to_id"`
	RepeatedJob             int       `json:"repeated_system_number"`
}

type WorkOrderBookingRequest struct {
	// Basic information
	BatchSystemNumber       int       `json:"batch_system_number"`
	BookingSystemNumber     int       `json:"booking_system_number"`
	EstimationSystemNumber  int       `json:"estimation_system_number"`
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	WorkOrderTypeName       string    `json:"work_order_type_name"`
	ServiceAdvisorId        int       `json:"service_advisor_id"`
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VariantId               int       `json:"variant_id"`
	ServiceSite             string    `json:"servicesite"`
	VehicleId               int       `json:"vehicle_id"`
	CustomerId              int       `json:"customer_id"`
	BilltoCustomerId        int       `json:"billto_customer_id"`
	CampaignId              int       `json:"campaign_id"`
	AgreementId             int       `json:"agreement_id"`
	ContractSystemNumber    int       `json:"contract_system_number"`
	CompanyId               int       `json:"company_id"`
	DealerRepresentativeId  int       `json:"dealer_representative_id"`
	FromEra                 bool      `json:"from_era"`
	QueueSystemNumber       int       `json:"queue_system_number"`
	WorkOrderArrivalTime    time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage int       `json:"work_order_current_mileage"`
	Storing                 bool      `json:"storing"`
	WorkOrderRemark         string    `json:"work_order_remark"`
	Unregistered            bool      `json:"unregistered"`
	WorkOrderProfitCenterId int       `json:"work_order_profit_center"`

	// Work order status and details
	WorkOrderEraNo           string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  time.Time `json:"work_order_era_expired_date"`
	WorkOrderStatusId        int       `json:"work_order_status_id"`
	WorkOrderStatusInsurance bool      `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsuranceCheck       bool      `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64   `json:"estimation_duration"`
	CustomerExpress    bool      `json:"customer_express"`
	LeaveCar           bool      `json:"leave_car"`
	CarWash            bool      `json:"car_wash"`
	PromiseDate        time.Time `json:"promise_date"`
	PromiseTime        time.Time `json:"promise_time"`

	// Additional information
	FSCouponNo        string  `json:"fs_coupon_no"`
	Notes             string  `json:"notes"`
	Suggestion        string  `json:"suggestion"`
	DownpaymentAmount float64 `json:"dp_amount"`
}

type WorkOrderBooking struct {

	// Basic information
	WorkOrderSystemNumber      int    `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderDocumentNumber    string `json:"work_order_document_number" parent_entity:"trx_work_order"`
	EstimationSystemNumber     int    `json:"estimation_system_number" parent_entity:"trx_work_order"`
	BookingSystemNumber        int    `json:"booking_system_number" parent_entity:"trx_work_order"`
	ServiceRequestSystemNumber int    `json:"service_request_system_number" parent_entity:"trx_work_order"`
	WorkOrderTypeId            int    `json:"work_order_type_id" parent_entity:"trx_work_order"`
	ServiceAdvisorId           int    `json:"service_advisor_id" parent_entity:"trx_work_order"`
	BrandId                    int    `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId                    int    `json:"model_id" parent_entity:"trx_work_order"`
	VariantId                  int    `json:"variant_id" parent_entity:"trx_work_order"`
	VehicleId                  int    `json:"vehicle_id" parent_entity:"trx_work_order"`
	CustomerId                 int    `json:"customer_id" parent_entity:"trx_work_order"`
	BilltoCustomerId           int    `json:"billable_to_id" parent_entity:"trx_work_order"`
	StatusId                   int    `json:"work_order_status_id" parent_entity:"trx_work_order"`
}

type WorkOrderBookingResponse struct {
	WorkOrderSystemNumber         int                             `json:"work_order_system_number"`
	WorkOrderDocumentNumber       string                          `json:"work_order_document_number"`
	WorkOrderDate                 string                          `json:"work_order_date"`
	EstimationSystemNumber        int                             `json:"estimation_system_number"`
	EstimationDocumentNumber      string                          `json:"estimation_document_number"`
	BookingSystemNumber           int                             `json:"booking_system_number"`
	BookingDocumentNumber         string                          `json:"booking_document_number"`
	ServiceRequestSystemNumber    int                             `json:"service_request_system_number"`
	WorkOrderTypeId               int                             `json:"work_order_type_id"`
	WorkOrderTypeName             string                          `json:"work_order_type_name"`
	WorkOrderStatusId             int                             `json:"work_order_status_id"`
	WorkOrderStatusName           string                          `json:"work_order_status_description"`
	ServiceAdvisorId              int                             `json:"service_advisor_id"`
	BrandId                       int                             `json:"brand_id"`
	BrandName                     string                          `json:"brand_name"`
	ModelId                       int                             `json:"model_id"`
	ModelName                     string                          `json:"model_name"`
	VariantId                     int                             `json:"variant_id"`
	VariantName                   string                          `json:"variant_name"`
	ServiceSite                   string                          `json:"service_site"`
	VehicleId                     int                             `json:"vehicle_id"`
	VehicleCode                   string                          `json:"vehicle_chassis_number"`
	VehicleTnkb                   string                          `json:"vehicle_registration_certificate_tnkb"`
	CustomerId                    int                             `json:"customer_id"`
	BilltoCustomerId              int                             `json:"billable_to_id"`
	CampaignId                    int                             `json:"campaign_id"`
	AgreementId                   int                             `json:"agreement_id"`
	ContractSystemNumber          int                             `json:"contract_system_number"`
	FromEra                       bool                            `json:"from_era"`
	QueueSystemNumber             int                             `json:"queue_system_number"`
	WorkOrderArrivalTime          time.Time                       `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage       int                             `json:"work_order_current_mileage"`
	WorkOrderRemark               string                          `json:"work_order_remark"`
	DealerRepresentativeId        int                             `json:"dealer_representative_id"`
	CompanyId                     int                             `json:"company_id"`
	WorkOrderProfitCenterId       int                             `json:"work_order_profit_center"`
	Titleprefix                   string                          `json:"title_prefix"`
	NameCust                      string                          `json:"name_customer"`
	PhoneCust                     string                          `json:"phone_customer"`
	MobileCust                    string                          `json:"mobile_customer"`
	MobileCustAlternative         string                          `json:"mobile_customer_alternative"`
	MobileCustDriver              string                          `json:"mobile_customer_driver"`
	ContactVia                    string                          `json:"contact_via"`
	WorkOrderInsuranceCheck       bool                            `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string                          `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time                       `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string                          `json:"insurance_claim_no"`
	WorkOrderEraExpiredDate       time.Time                       `json:"work_order_era_expired_date"`
	PromiseDate                   time.Time                       `json:"promise_date"`
	PromiseTime                   time.Time                       `json:"promise_time"`
	EstimationDuration            float64                         `json:"estimation_duration"`
	WorkOrderInsuranceOwnRisk     float64                         `json:"insurance_own_risk"`
	CustomerExpress               bool                            `json:"customer_express"`
	LeaveCar                      bool                            `json:"leave_car"`
	CarWash                       bool                            `json:"car_wash"`
	FSCouponNo                    string                          `json:"fs_coupon_no"`
	Notes                         string                          `json:"notes"`
	Suggestion                    string                          `json:"suggestion"`
	DownpaymentAmount             float64                         `json:"dp_amount"`
	WorkOrderInsurancePic         string                          `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string                          `json:"insurance_workorder_number"`
	WorkOrderEraNo                string                          `json:"work_order_era_no"`
	WorkOrderStatusInsurance      bool                            `json:"work_order_status_insurance"`
	WorkOrderDetailService        WorkOrderDetailsResponseRequest `json:"work_order_details_service"`
	WorkOrderDetailVehicle        WorkOrderDetailsResponseVehicle `json:"work_order_details_vehicle"`
	WorkOrderDetails              WorkOrderDetailsResponse        `json:"work_order_details"`
}

type WorkOrderAffiliateRequest struct {
	// Basic information
	ServiceRequestId        int       `json:"service_request_id"`
	ServiceRequestNumber    int       `json:"service_request_number"`
	ServiceRequestDate      time.Time `json:"service_request_date"`
	ServiceRequestCompany   string    `json:"service_request_company"`
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	ServiceAdvisorId        int       `json:"service_advisor_id"`
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VariantId               int       `json:"variant_id"`
	ServiceSite             string    `json:"servicesite"`
	VehicleId               int       `json:"vehicle_id"`
	CustomerId              int       `json:"customer_id"`
	BilltoCustomerId        int       `json:"billto_customer_id"`
	CampaignId              int       `json:"campaign_id"`
	AgreementId             int       `json:"agreement_id"`
	BoookingId              int       `json:"booking_system_number"`
	EstimationId            int       `json:"estimation_system_number"`
	ContractSystemNumber    int       `json:"contract_system_number"`
	CompanyId               int       `json:"company_id"`
	DealerRepresentativeId  int       `json:"dealer_representative_id"`
	FromEra                 bool      `json:"from_era"`
	QueueSystemNumber       int       `json:"queue_system_number"`
	WorkOrderArrivalTime    time.Time `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage int       `json:"work_order_current_mileage"`
	Storing                 bool      `json:"storing"`
	WorkOrderRemark         string    `json:"work_order_remark"`
	Unregistered            bool      `json:"unregistered"`
	WorkOrderProfitCenter   int       `json:"work_order_profit_center"`

	// Work order status and details
	WorkOrderEraNo           string    `json:"work_order_era_no"`
	WorkOrderEraExpiredDate  time.Time `json:"work_order_era_expired_date"`
	WorkOrderStatusId        int       `json:"work_order_status_id"`
	WorkOrderStatusInsurance bool      `json:"work_order_status_insurance"`

	// Customer contact information
	Titleprefix           string `json:"title_prefix"`
	NameCust              string `json:"name_customer"`
	PhoneCust             string `json:"phone_customer"`
	MobileCust            string `json:"mobile_customer"`
	MobileCustAlternative string `json:"mobile_customer_alternative"`
	MobileCustDriver      string `json:"mobile_customer_driver"`
	ContactVia            string `json:"contact_via"`

	// Insurance details
	WorkOrderInsurancePolicyNo    string    `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float64   `json:"insurance_own_risk"`

	// Estimation and service details
	EstimationDuration float64   `json:"estimation_duration"`
	CustomerExpress    bool      `json:"customer_express"`
	LeaveCar           bool      `json:"leave_car"`
	CarWash            bool      `json:"car_wash"`
	PromiseDate        time.Time `json:"promise_date"`
	PromiseTime        time.Time `json:"promise_time"`

	// Additional information
	FSCouponNo        string  `json:"fs_coupon_no"`
	Notes             string  `json:"notes"`
	Suggestion        string  `json:"suggestion"`
	DownpaymentAmount float64 `json:"dp_amount"`
}

type WorkOrderAffiliateGetResponse struct {
	WorkOrderSystemNumber        int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber      string `json:"work_order_document_number"`
	ServiceRequestSystemNumber   int    `json:"service_request_system_number"`
	ServiceRequestDate           string `json:"service_request_date"`
	ServiceRequestDocumentNumber string `json:"service_request_document_number"`
	BrandId                      int    `json:"brand_id"`
	BrandName                    string `json:"brand_name"`
	ModelId                      int    `json:"model_id"`
	ModelName                    string `json:"model_name"`
	VehicleId                    int    `json:"vehicle_id"`
	VehicleCode                  string `json:"vehicle_chassis_number"`
	VehicleTnkb                  string `json:"vehicle_registration_certificate_tnkb"`
	CompanyId                    int    `json:"company_id"`
	CompanyName                  string `json:"company_name"`
}

type WorkOrderAffiliateResponse struct {
	WorkOrderSystemNumber         int                             `json:"work_order_system_number"`
	WorkOrderDocumentNumber       string                          `json:"work_order_document_number"`
	WorkOrderDate                 string                          `json:"work_order_date"`
	EstimationSystemNumber        int                             `json:"estimation_system_number"`
	EstimationDocumentNumber      string                          `json:"estimation_document_number"`
	BookingSystemNumber           int                             `json:"booking_system_number"`
	BookingDocumentNumber         string                          `json:"booking_document_number"`
	ServiceRequestSystemNumber    int                             `json:"service_request_system_number"`
	ServiceRequestDocumentNumber  string                          `json:"service_request_document_number"`
	WorkOrderTypeId               int                             `json:"work_order_type_id"`
	WorkOrderTypeName             string                          `json:"work_order_type_name"`
	WorkOrderStatusId             int                             `json:"work_order_status_id"`
	WorkOrderStatusName           string                          `json:"work_order_status_description"`
	ServiceAdvisorId              int                             `json:"service_advisor_id"`
	BrandId                       int                             `json:"brand_id"`
	BrandName                     string                          `json:"brand_name"`
	ModelId                       int                             `json:"model_id"`
	ModelName                     string                          `json:"model_name"`
	VariantId                     int                             `json:"variant_id"`
	VariantName                   string                          `json:"variant_name"`
	ServiceSite                   string                          `json:"service_site"`
	VehicleId                     int                             `json:"vehicle_id"`
	VehicleCode                   string                          `json:"vehicle_chassis_number"`
	VehicleTnkb                   string                          `json:"vehicle_registration_certificate_tnkb"`
	CustomerId                    int                             `json:"customer_id"`
	BilltoCustomerId              int                             `json:"billable_to_id"`
	CampaignId                    int                             `json:"campaign_id"`
	AgreementId                   int                             `json:"agreement_id"`
	ContractSystemNumber          int                             `json:"contract_system_number"`
	FromEra                       bool                            `json:"from_era"`
	QueueSystemNumber             int                             `json:"queue_system_number"`
	WorkOrderArrivalTime          time.Time                       `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage       int                             `json:"work_order_current_mileage"`
	WorkOrderRemark               string                          `json:"work_order_remark"`
	DealerRepresentativeId        int                             `json:"dealer_representative_id"`
	CompanyId                     int                             `json:"company_id"`
	WorkOrderProfitCenterId       int                             `json:"work_order_profit_center"`
	Titleprefix                   string                          `json:"title_prefix"`
	NameCust                      string                          `json:"name_customer"`
	PhoneCust                     string                          `json:"phone_customer"`
	MobileCust                    string                          `json:"mobile_customer"`
	MobileCustAlternative         string                          `json:"mobile_customer_alternative"`
	MobileCustDriver              string                          `json:"mobile_customer_driver"`
	ContactVia                    string                          `json:"contact_via"`
	WorkOrderInsuranceCheck       bool                            `json:"insurance_check"`
	WorkOrderInsurancePolicyNo    string                          `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time                       `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string                          `json:"insurance_claim_no"`
	WorkOrderEraExpiredDate       time.Time                       `json:"work_order_era_expired_date"`
	PromiseDate                   time.Time                       `json:"promise_date"`
	PromiseTime                   time.Time                       `json:"promise_time"`
	EstimationDuration            float64                         `json:"estimation_duration"`
	WorkOrderInsuranceOwnRisk     float64                         `json:"insurance_own_risk"`
	CustomerExpress               bool                            `json:"customer_express"`
	LeaveCar                      bool                            `json:"leave_car"`
	CarWash                       bool                            `json:"car_wash"`
	FSCouponNo                    string                          `json:"fs_coupon_no"`
	Notes                         string                          `json:"notes"`
	Suggestion                    string                          `json:"suggestion"`
	DownpaymentAmount             float64                         `json:"dp_amount"`
	WorkOrderInsurancePic         string                          `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string                          `json:"insurance_workorder_number"`
	WorkOrderEraNo                string                          `json:"work_order_era_no"`
	WorkOrderStatusInsurance      bool                            `json:"work_order_status_insurance"`
	WorkOrderDetailService        WorkOrderDetailsResponseRequest `json:"work_order_details_service"`
	WorkOrderDetailVehicle        WorkOrderDetailsResponseVehicle `json:"work_order_details_vehicle"`
	WorkOrderDetails              WorkOrderDetailsResponse        `json:"work_order_details"`
}

type WorkOrderLookupRequest struct {
	BrandId int `json:"brand_id"`
	ModelId int `json:"model_id"`
}

type WorkOrderLookupResponse struct {
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VehicleId               int       `json:"vehicle_id"`
	CustomerId              int       `json:"customer_id"`
	VehicleCode             string    `json:"vehicle_chassis_number"`
	VehicleTnkb             string    `json:"vehicle_registration_certificate_tnkb"`
	VehicleCertificateOwner string    `json:"vehicle_registration_certificate_owner_name"`
	VehicleProduction       string    `json:"vehicle_production_year"`
	VehicleVariantColour    string    `json:"variant_colour_description"`
	VehicleServiceBookingNo string    `json:"service_booking_number"`
	VehicleLastServiceDate  time.Time `json:"last_service_date"`
	VehicleLastKm           int       `json:"last_km"`
}

type WorkOrderVehicleResponse struct {
	VehicleId               int       `json:"vehicle_id"`
	VehicleCode             string    `json:"vehicle_chassis_number"`
	VehicleTnkb             string    `json:"vehicle_registration_certificate_tnkb"`
	VehicleCertificateOwner string    `json:"vehicle_registration_certificate_owner_name"`
	VehicleProduction       string    `json:"vehicle_production_year"`
	VehicleVariantColour    string    `json:"variant_colour_description"`
	VehicleServiceBookingNo string    `json:"service_booking_number"`
	VehicleLastServiceDate  time.Time `json:"last_service_date"`
	VehicleLastKm           int       `json:"last_km"`
}

type CustomerResponse struct {
	CustomerId   int    `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	CustomerType string `json:"customer_type"`
	CustomerCode string `json:"customer_code"`
}

type WorkOrderCampaignResponse struct {
	CampaignId         int       `json:"campaign_id"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	CampaignPeriodFrom time.Time `json:"campign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
}

type WorkOrderBillable struct {
	BillableToName string `json:"billable_to_name"`
	BillableToID   int    `json:"billable_to_id"`
	IsActive       bool   `json:"is_active"`
	BillableToCode string `json:"billable_to_code"`
}
type WorkOrderDropPoint struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}
type WorkOrderVehicleBrand struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type WorkOrderVehicleModel struct {
	ModelId              int    `json:"model_id"`
	ModelCode            string `json:"model_code"`
	ModelName            string `json:"model_description"`
	ModelCodeDescription string `json:"model_code_description"`
}

type WorkOrderVehicleVariant struct {
	VariantId   int    `json:"variant_id"`
	VariantCode string `json:"variant_code"`
	VariantName string `json:"variant_description"`
}

type WorkOrderVehicleColour struct {
	VariantColourId   int    `json:"colour_id"`
	VariantColourCode string `json:"colour_commercial_name"`
	VariantColourName string `json:"colour_police_name"`
}

type WorkOrderTypeResponse struct {
	WorkOrderTypeId   int    `json:"work_order_type_id"`
	WorkOrderTypeCode string `json:"work_order_type_code"`
	WorkOrderTypeName string `json:"work_order_type_description"`
}

type WorkOrderStatusResponse struct {
	WorkOrderStatusId   int    `json:"work_order_status_id"`
	WorkOrderStatusCode string `json:"work_order_status_code"`
	WorkOrderStatusName string `json:"work_order_status_description"`
}

type WorkOrderServiceRequest struct {
	WorkOrderServiceId     int    `json:"work_order_service_id"`
	WorkOrderSystemNumber  int    `json:"work_order_system_number"`
	WorkOrderServiceRemark string `json:"work_order_service_remark"`
}

type WorkOrderServiceResponse struct {
	WorkOrderServiceId     int    `json:"work_order_service_id"`
	WorkOrderSystemNumber  int    `json:"work_order_system_number"`
	WorkOrderServiceRemark string `json:"work_order_service_remark"`
}

type WorkOrderServiceVehicleRequest struct {
	WorkOrderSystemNumber  int       `json:"work_order_system_number"`
	WorkOrderVehicleDate   time.Time `json:"work_order_vehicle_date"`
	WorkOrderVehicleRemark string    `json:"work_order_vehicle_remark"`
}

type WorkOrderServiceVehicleResponse struct {
	WorkOrderSystemNumber  int       `json:"work_order_system_number"`
	WorkOrderVehicleDate   time.Time `json:"work_order_vehicle_date"`
	WorkOrderVehicleRemark string    `json:"work_order_vehicle_remark"`
}

type WorkOrderAddRequest struct {
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	WorkOrderStatusId       int       `json:"work_order_status_id"`
}

type WorkOrderDetailRequest struct {
	WorkOrderDetailId     int     `json:"work_order_detail_id" parent_entity:"trx_work_order_detail" main_table:"trx_work_order_detail"`
	WorkOrderSystemNumber int     `json:"work_order_system_number" parent_entity:"trx_work_order_detail"`
	LineTypeId            int     `json:"line_type_id" parent_entity:"trx_work_order_detail"`
	TransactionTypeId     int     `json:"transaction_type_id" parent_entity:"trx_work_order_detail" `
	JobTypeId             int     `json:"job_type_id" parent_entity:"trx_work_order_detail"`
	FrtQuantity           float64 `json:"frt_quantity" parent_entity:"trx_work_order_detail"`
	SupplyQuantity        float64 `json:"supply_quantity" parent_entity:"trx_work_order_detail"`
	PriceListId           int     `json:"price_list_id" parent_entity:"trx_work_order_detail"`
	WarehouseId           int     `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	ItemId                int     `json:"item_id" parent_entity:"trx_work_order_detail"`
	OperationId           int     `json:"operation_id" parent_entity:"trx_work_order_detail"`
	OperationItemCode     string  `json:"operation_item_code" parent_entity:"trx_work_order_detail"`
	ProposedPrice         float64 `json:"operation_item_discount_request_amount" parent_entity:"trx_work_order_detail"`
	OperationItemPrice    float64 `json:"operation_item_price" parent_entity:"trx_work_order_detail"`
}

type WorkOrderDetailResponse struct {
	WorkOrderDetailId                  int     `json:"work_order_detail_id"`
	WorkOrderSystemNumber              int     `json:"work_order_system_number"`
	LineTypeId                         int     `json:"line_type_id"`
	TransactionTypeId                  int     `json:"transaction_type_id"`
	JobTypeId                          int     `json:"job_type_id"`
	WarehouseId                        int     `json:"warehouse_id"`
	ItemId                             int     `json:"item_id"`
	FrtQuantity                        float64 `json:"frt_quantity"`
	SupplyQuantity                     float64 `json:"supply_quantity"`
	OperationItemPrice                 float64 `json:"operation_item_price"`
	OperationItemDiscountAmount        float64 `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount float64 `json:"operation_item_discount_request_amount"`
}

type WorkOrderAffiliate struct {
	WorkOrderSystemNumber      int    `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderDocumentNumber    string `json:"work_order_document_number" parent_entity:"trx_work_order"`
	ServiceRequestSystemNumber int    `json:"service_request_system_number" parent_entity:"trx_work_order"`
	BrandId                    int    `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId                    int    `json:"model_id" parent_entity:"trx_work_order"`
	VehicleId                  int    `json:"vehicle_id" parent_entity:"trx_work_order"`
	CompanyId                  int    `json:"company_id" parent_entity:"trx_work_order"`
}

type WorkOrderAffiliatedRequest struct {
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	ServiceRequestId        int       `json:"service_request_id"`
	ServiceRequestNumber    int       `json:"service_request_number"`
	ServiceRequestDate      time.Time `json:"service_request_date"`
	ServiceRequestCompany   string    `json:"service_request_company"`
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VehicleId               int       `json:"vehicle_id"`
}

type WorkOrderStatusRequest struct {
	WorkOrderStatusId   int    `json:"work_order_status_id"`
	WorkOrderStatusCode string `json:"work_order_status_code"`
	WorkOrderStatusName string `json:"work_order_status_description"`
}

type WorkOrderTypeRequest struct {
	WorkOrderTypeId   int    `json:"work_order_type_id"`
	WorkOrderTypeCode string `json:"work_order_type_code"`
	WorkOrderTypeName string `json:"work_order_type_description"`
}

type WorkOrderBillableRequest struct {
	BillableToName string `json:"billable_to_name"`
	BillableToID   int    `json:"billable_to_id"`
	IsActive       bool   `json:"is_active"`
	BillableToCode string `json:"billable_to_code"`
}

type WorkOrderDropPointRequest struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}

type SubmitWorkOrderResponse struct {
	DocumentNumber        string `json:"document_number"`
	WorkOrderSystemNumber int    `json:"work_order_system_number"`
}

type ChangeBillToRequest struct {
	WorkOrderSystemNumber int `json:"work_order_system_number"`
	BillToCustomerId      int `json:"customer_id"`
}

type ChangePhoneNoRequest struct {
	WorkOrderSystemNumber int    `json:"work_order_system_number"`
	PhoneNo               string `json:"contact_person_phone"`
}

type BrandDocResponse struct {
	BrandId           int    `json:"brand_id"`
	BrandCode         string `json:"brand_code"`
	BrandName         string `json:"brand_name"`
	BrandAbbreviation string `json:"brand_abbreveation"`
}

type VehicleResponse struct {
	VehicleId           int             `json:"vehicle_id"`
	VehicleCode         string          `json:"vehicle_chassis_number"`
	VehicleTnkb         string          `json:"vehicle_registration_certificate_tnkb"`
	VehicleProduction   json.RawMessage `json:"vehicle_production_year"`
	VehicleLastKm       json.RawMessage `json:"vehicle_last_km"`
	VehicleBrandId      int             `json:"vehicle_brand_id"`
	VehicleModelId      int             `json:"vehicle_model_id"`
	VehicleModelVariant string          `json:"model_variant_colour_description"`
	VehicleVariantId    int             `json:"vehicle_variant_id"`
	VehicleColourId     int             `json:"vehicle_colour_id"`
	VehicleOwner        string          `json:"vehicle_registration_certificate_owner_name"`
}

type Linetype struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
}

type WorkOrderTransactionType struct {
	TransactionTypeId   int    `json:"work_order_transaction_type_id"`
	TransactionTypeCode string `json:"work_order_transaction_type_code"`
	TransactionTypeName string `json:"work_order_transaction_type_description"`
}
