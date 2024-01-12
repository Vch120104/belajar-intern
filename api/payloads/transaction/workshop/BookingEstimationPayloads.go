package transactionworkshoppayloads

import "time"

type SaveBookingEstimationRequest struct {
	BatchSystemNumber              int       `gorm:"column:batch_system_number;not null;primaryKey" json:"batch_system_number"`
	BookingSystemNumber            int       `json:"booking_system_number"`
	BrandId                        int       `json:"brand_id"`
	ModelId                        int       `json:"model_id"`
	VariantId                      int       `json:"variant_id"`
	VehicleId                      int       `json:"vehicle_id"`
	EstimationSystemNumber         int       `json:"estimation_system_number"`
	PdiSystemNumber                int       `json:"pdi_system_number"`
	ServiceRequestSystemNumber     int       `json:"system_request_system_number"`
	ContractSystemNumber           int       `json:"contract_system_number"`
	AgreementId                    int       `json:"agreement_id"`
	CampaignId                     int       `json:"campaign_id"`
	CompanyId                      int       `json:"company_id"`
	ProfitCenterId                 int       `json:"profit_center_id"`
	DealerRepresentativeId         int       `json:"dealer_representative_id"`
	CustomerId                     int       `json:"customer_id"`
	DocumentStatusId               int       `json:"document_status_id"`
	BookingEstimationBatchDate     time.Time `json:"booking_estimation_batch_date"`
	BookingEstimationVehicleNumber string    `json:"booking_estimation_vehicle_number"`
	AgreementNumberBr              string    `json:"agreement_number_br"`
	IsUnregistered                 string    `json:"is_unregistered"`
	ContactPersonName              string    `json:"contact_person_name"`
	ContactPersonPhone             string    `json:"contact_person_phone"`
	ContactPersonMobile            string    `json:"contact_person_mobile"`
	ContactPersonVia               string    `json:"contact_person_via"`
	InsurancePolicyNo              string    `json:"insurance_policy_no"`
	InsuranceExpiredDate           time.Time `json:"insurance_expired_date"`
	InsuranceClaimNo               string    `json:"insurance_claim_no"`
	InsurancePic                   string    `json:"insurance_pic"`
}
