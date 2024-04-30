package transactionworkshoppayloads

import "time"

type WorkOrderResponse struct {
	BatchSystemNumber             int           `json:"batch_system_number"`
	WorkOrderSystemNumber         string        `json:"work_order_system_number"`
	WorkOrderDate                 time.Time     `json:"work_order_date"`
	WorkOrderTypeId               int           `json:"work_order_type_id"`
	WorkOrderServiceSite          int           `json:"work_order_servicesite_id"`
	BrandId                       int           `json:"brand_id"`
	ModelId                       int           `json:"model_id"`
	VariantId                     int           `json:"variant_id"`
	VehicleId                     int           `json:"vehicle_id"`
	BilltoCustomerId              int           `json:"customer_id"`
	WorkOrderStatusEra            bool          `json:"work_order_status_era"`
	WorkOrderEraNo                string        `json:"work_order_era_no"`
	WorkOrderEraExpiredDate       time.Time     `json:"work_order_era_expired_date"`
	QueueSystemNumber             int           `json:"queue_system_number"`
	WorkOrderArrivalTime          time.Time     `json:"work_order_arrival_time"`
	WorkOrderCurrentMileage       int           `json:"work_order_current_mileage"`
	WorkOrderStatusStoring        bool          `json:"work_order_status_storing"`
	WorkOrderRemark               string        `json:"work_order_remark"`
	WorkOrderStatusUnregistered   bool          `json:"work_order_status_unregistered"`
	WorkOrderProfitCenter         string        `json:"work_order_profit_center"`
	WorkOrderDealerRepCode        string        `json:"work_order_dealer_rep_code"`
	CampaignId                    int           `json:"campaign_id"`
	AgreementId                   int           `json:"agreement_id"`
	ServiceRequestSystemNumber    int           `json:"system_request_system_number"`
	EstimationSystemNumber        int           `json:"estimation_system_number"`
	ContractSystemNumber          int           `json:"contract_system_number"`
	CompanyId                     int           `json:"company_id"`
	DealerRepresentativeId        int           `json:"dealer_representative_id"`
	Titleprefix                   string        `json:"title_prefix"`
	NameCust                      string        `json:"name_customer"`
	PhoneCust                     string        `json:"phone_customer"`
	MobileCust                    string        `json:"mobile_customer"`
	MobileCustAlternative         string        `json:"mobile_customer_alternative"`
	MobileCustDriver              string        `json:"mobile_customer_driver"`
	ContactVia                    string        `json:"contact_via"`
	WorkOrderStatusInsurance      bool          `json:"work_order_status_insurance"`
	WorkOrderInsurancePolicyNo    string        `json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time     `json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string        `json:"insurance_claim_no"`
	WorkOrderInsurancePic         string        `json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string        `json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float32       `json:"insurance_own_risk"`
	BrandDetails                  BrandResponse `json:"brand_details"`
}

type WorkOrderRequest struct {
	BatchSystemNumber             int       `json:"batch_system_number"`
	WorkOrderSystemNumber         string    `json:"work_order_system_number"`
	WorkOrderDate                 time.Time `json:"work_order_date"`
	WorkOrderTypeId               int       `json:"work_order_type_id"`
	WorkOrderServiceSite          int       `json:"work_order_servicesite_id"`
	BrandId                       int       `json:"brand_id"`
	ModelId                       int       `json:"model_id"`
	VariantId                     int       `json:"variant_id"`
	VehicleId                     int       `json:"vehicle_id"`
	BilltoCustomerId              int       `json:"customer_id"`
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

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}
