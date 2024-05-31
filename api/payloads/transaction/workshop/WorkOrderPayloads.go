package transactionworkshoppayloads

import "time"

type WorkOrderResponse struct {
	WorkOrderStatusId             int       `json:"work_order_status_id"`
	WorkOrderSystemNumber         int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber       string    `json:"work_order_document_number"`
	WorkOrderDate                 time.Time `json:"work_order_date"`
	WorkOrderTypeId               int       `json:"work_order_type_id"`
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
	WorkOrderInsuranceOwnRisk     float32   `json:"insurance_own_risk"`
}

type WorkOrderRequest struct {
	WorkOrderStatusId             int       `json:"work_order_status_id"`
	WorkOrderSystemNumber         int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber       string    `json:"work_order_document_number"`
	WorkOrderDate                 time.Time `json:"work_order_date"`
	WorkOrderTypeId               int       `json:"work_order_type_id"`
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
	WorkOrderInsuranceOwnRisk     float32   `json:"insurance_own_risk"`
}

type WorkOrderLookupRequest struct {
	WorkOrderSystemNumber   string `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderDocumentNumber string `json:"work_order_document_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	VehicleId               int    `json:"vehicle_id" parent_entity:"trx_work_order"`
	CustomerId              int    `json:"customer_id" parent_entity:"trx_work_order"`
}

type WorkOrderLookupResponse struct {
	WorkOrderDocumentNumber int    `json:"work_order_document_number"`
	WorkOrderSystemNumber   string `json:"work_order_system_number"`
	VehicleId               int    `json:"vehicle_id"`
	CustomerId              int    `json:"customer_id"`
}

type WorkOrderVehicleResponse struct {
	VehicleId               int       `json:"vehicle_id"`
	VehicleCode             string    `json:"vehicle_chassis_number"`
	VehicleTnkb             string    `json:"registration_certificate_tnkb"`
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

type WorkOrderServiceRequest struct {
	WorkOrderServiceId      int    `json:"work_order_service_id"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderServiceRemark  string `json:"work_order_service_remark"`
}

type WorkOrderServiceVehicleRequest struct {
	WorkOrderVehicleId      int       `json:"work_order_vehicle_id"`
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderVehicleDate    time.Time `json:"work_order_vehicle_date"`
	WorkOrderVehicleRemark  string    `json:"work_order_vehicle_remark"`
}

type WorkOrderAddRequest struct {
	WorkOrderSystemNumber   int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	WorkOrderTypeId         int       `json:"work_order_type_id"`
	WorkOrderStatusId       int       `json:"work_order_status_id"`
}
