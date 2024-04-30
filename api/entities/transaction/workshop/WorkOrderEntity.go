package transactionworkshopentities

import "time"

const TableNameWorkOrder = "trx_work_order"

type WorkOrder struct {
	BatchSystemNumber             int       `gorm:"column:batch_system_number;size:30;" json:"batch_system_number"`
	WorkOrderSystemNumber         string    `gorm:"column:work_order_system_number;type:varchar(18);primaryKey" json:"work_order_system_number"`
	WorkOrderDate                 time.Time `gorm:"column:work_order_date" json:"work_order_date"`
	WorkOrderTypeId               int       `gorm:"column:work_order_type_id;size:30;" json:"work_order_type_id"`
	WorkOrderServiceSite          int       `gorm:"column:work_order_servicesite_id;size:30;" json:"work_order_servicesite_id"`
	BrandId                       int       `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId                       int       `gorm:"column:model_id;size:30;" json:"model_id"`
	VariantId                     int       `gorm:"column:variant_id;size:30;" json:"variant_id"`
	VehicleId                     int       `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	BilltoCustomerId              int       `gorm:"column:customer_id" json:"customer_id"`
	WorkOrderStatusEra            bool      `gorm:"column:work_order_status_era" json:"work_order_status_era"`
	WorkOrderEraNo                string    `gorm:"column:work_order_era_no;type:varchar(25)" json:"work_order_era_no"`
	WorkOrderEraExpiredDate       time.Time `gorm:"column:work_order_era_expired_date" json:"work_order_era_expired_date"`
	QueueSystemNumber             int       `gorm:"column:queue_system_number" json:"queue_system_number"`
	WorkOrderArrivalTime          time.Time `gorm:"column:work_order_arrival_time" json:"work_order_arrival_time"`
	WorkOrderCurrentMileage       int       `gorm:"column:work_order_current_mileage" json:"work_order_current_mileage"`
	WorkOrderStatusStoring        bool      `gorm:"column:work_order_status_storing" json:"work_order_status_storing"`
	WorkOrderRemark               string    `gorm:"column:work_order_remark;type:varchar(200)" json:"work_order_remark"`
	WorkOrderStatusUnregistered   bool      `gorm:"column:work_order_status_unregistered" json:"work_order_status_unregistered"`
	WorkOrderProfitCenter         string    `gorm:"column:work_order_profit_center;type:varchar(200)" json:"work_order_profit_center"`
	WorkOrderDealerRepCode        string    `gorm:"column:work_order_dealer_rep_code;type:varchar(200)" json:"work_order_dealer_rep_code"`
	CampaignId                    int       `gorm:"column:campaign_id" json:"campaign_id"`
	AgreementId                   int       `gorm:"column:agreement_id" json:"agreement_id"`
	ServiceRequestSystemNumber    int       `gorm:"column:service_request_system_number" json:"system_request_system_number"`
	EstimationSystemNumber        int       `gorm:"column:estimation_system_number" json:"estimation_system_number"`
	ContractSystemNumber          int       `gorm:"column:contract_system_number" json:"contract_system_number"`
	CompanyId                     int       `gorm:"column:company_id" json:"company_id"`
	DealerRepresentativeId        int       `gorm:"column:dealer_representative_id" json:"dealer_representative_id"`
	Titleprefix                   string    `gorm:"column:title_prefix" json:"title_prefix"`
	NameCust                      string    `gorm:"column:name_customer" json:"name_customer"`
	PhoneCust                     string    `gorm:"column:phone_customer" json:"phone_customer"`
	MobileCust                    string    `gorm:"column:mobile_customer" json:"mobile_customer"`
	MobileCustAlternative         string    `gorm:"column:mobile_customer_alternative" json:"mobile_customer_alternative"`
	MobileCustDriver              string    `gorm:"column:mobile_customer_driver" json:"mobile_customer_driver"`
	ContactVia                    string    `gorm:"column:contact_via" json:"contact_via"`
	WorkOrderStatusInsurance      bool      `gorm:"column:work_order_status_insurance" json:"work_order_status_insurance"`
	WorkOrderInsurancePolicyNo    string    `gorm:"column:insurance_policy_no;type:varchar(25)" json:"insurance_policy_no"`
	WorkOrderInsuranceExpiredDate time.Time `gorm:"column:insurance_expired_date" json:"insurance_expired_date"`
	WorkOrderInsuranceClaimNo     string    `gorm:"column:insurance_claim_no;type:varchar(25)" json:"insurance_claim_no"`
	WorkOrderInsurancePic         string    `gorm:"column:insurance_pic;type:varchar(35)" json:"insurance_pic"`
	WorkOrderInsuranceWONumber    string    `gorm:"column:insurance_workorder_number;type:varchar(35)" json:"insurance_workorder_number"`
	WorkOrderInsuranceOwnRisk     float32   `gorm:"column:insurance_own_risk" json:"insurance_own_risk"`
}

func (*WorkOrder) TableName() string {
	return TableNameWorkOrder
}
