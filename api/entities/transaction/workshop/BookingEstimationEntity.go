package transactionworkshopentities

import "time"

const TableNameBookingEstimation = "trx_booking_estimation"

type BookingEstimation struct {
	BatchSystemNumber                int                                `gorm:"column:batch_system_number;not null;size:30;primaryKey" json:"batch_system_number"`
	BookingSystemNumber              int                                `gorm:"column:booking_system_number;size:30;" json:"booking_system_number"`
	BrandId                          int                                `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId                          int                                `gorm:"column:model_id;size:30;" json:"model_id"`
	VariantId                        int                                `gorm:"column:variant_id;size:30;" json:"variant_id"`
	VehicleId                        int                                `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	EstimationSystemNumber           int                                `gorm:"column:estimation_system_number;size:30;" json:"estimation_system_number"`
	PdiSystemNumber                  int                                `gorm:"column:pdi_system_number;size:30;" json:"pdi_system_number"`
	ServiceRequestSystemNumber       int                                `gorm:"column:service_request_system_number;size:30;" json:"system_request_system_number"`
	ContractSystemNumber             int                                `gorm:"column:contract_system_number;size:30;" json:"contract_system_number"`
	AgreementId                      int                                `gorm:"column:agreement_id;size:30;" json:"agreement_id"`
	CampaignId                       int                                `gorm:"column:campaign_id;size:30;" json:"campaign_id"`
	CompanyId                        int                                `gorm:"column:company_id;size:30;" json:"company_id"`
	ProfitCenterId                   int                                `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	DealerRepresentativeId           int                                `gorm:"column:dealer_representative_id;size:30;" json:"dealer_representative_id"`
	CustomerId                       int                                `gorm:"column:customer_id;size:30;" json:"customer_id"`
	DocumentStatusId                 int                                `gorm:"column:document_status_id;size:30;" json:"document_status_id"`
	BookingEstimationBatchDate       time.Time                          `gorm:"column:booking_estimation_batch_date" json:"booking_estimation_batch_date"`
	BookingEstimationVehicleNumber   string                             `gorm:"column:booking_estimation_vehicle_number;type:varchar(10)" json:"booking_estimation_vehicle_number"`
	AgreementNumberBr                string                             `gorm:"column:agreement_number_br;type:varchar(20)" json:"agreement_number_br"`
	IsUnregistered                   bool                               `gorm:"column:is_unregistered;" json:"is_unregistered"`
	ContactPersonName                string                             `gorm:"column:contact_person_name;type:varchar(40)" json:"contact_person_name"`
	ContactPersonPhone               string                             `gorm:"column:contact_person_phone;type:varchar(13)" json:"contact_person_phone"`
	ContactPersonMobile              string                             `gorm:"column:contact_person_mobile;type:varchar(13)" json:"contact_person_mobile"`
	ContactPersonViaId               int                                `gorm:"column:contact_person_via;type:varchar(5)" json:"contact_person_via"`
	InsurancePolicyNo                string                             `gorm:"column:insurance_policy_no;type:varchar(25)" json:"insurance_policy_no"`
	InsuranceExpiredDate             time.Time                          `gorm:"column:insurance_expired_date" json:"insurance_expired_date"`
	InsuranceClaimNo                 string                             `gorm:"column:insurance_claim_no;type:varchar(25)" json:"insurance_claim_no"`
	InsurancePic                     string                             `gorm:"column:insurance_pic;type:varchar(35)" json:"insurance_pic"`
	BookingEstimationAllocation      []BookingEstimationAllocation      `gorm:"foreignKey:BatchSystemNumber;references:BatchSystemNumber" json:"booking_estimation_allocation_batch"`
	BookingEstimationServiceDiscount []BookingEstimationServiceDiscount `gorm:"foreignKey:BatchSystemNumber;references:BatchSystemNumber" json:"booking_estimation_service_discount_batch"`
}

func (*BookingEstimation) TableName() string {
	return TableNameBookingEstimation
}
