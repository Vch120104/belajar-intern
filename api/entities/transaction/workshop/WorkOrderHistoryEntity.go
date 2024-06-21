package transactionworkshopentities

import "time"

var CreateWorkOrderHistoryTable = "trx_work_order_history"

type WorkOrderHistory struct {
	WorkOrderHistoryId                  int       `gorm:"column:work_order_history_id;size:30;not null;primaryKey" json:"work_order_history_id"`
	WorkOrderSystemNumber               int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	CompanyID                           int       `gorm:"column:company_id;size:30" json:"company_id"`
	WorkOrderDocumentNumber             string    `gorm:"column:work_order_document_number;size:100" json:"work_order_document_number"`
	WorkOrderStatusID                   string    `gorm:"column:work_order_status_id;size:100" json:"work_order_status_id"`
	WorkOrderDate                       time.Time `gorm:"column:work_order_date;type:datetime" json:"work_order_date"`
	WorkOrderCloseDate                  time.Time `gorm:"column:work_order_close_date;type:datetime" json:"work_order_close_date"`
	WorkOrderTypeID                     int       `gorm:"column:work_order_type_id;size:30" json:"work_order_type_id"`
	WorkOrderRepeatedSystemNumber       int       `gorm:"column:work_order_repeated_system_number;size:30" json:"work_order_repeated_system_number"`
	WorkOrderRepeatedDocumentNumber     string    `gorm:"column:work_order_repeated_document_number;size:100" json:"work_order_repeated_document_number"`
	ProfitCenterID                      int       `gorm:"column:profit_center_id;size:30" json:"profit_center_id"`
	BrandID                             int       `gorm:"column:brand_id;size:30" json:"brand_id"`
	ModelID                             int       `gorm:"column:model_id;size:30" json:"model_id"`
	VariantID                           int       `gorm:"column:variant_id;size:30" json:"variant_id"`
	VehicleID                           int       `gorm:"column:vehicle_id;size:30" json:"vehicle_id"`
	BillableToID                        int       `gorm:"column:billable_to_id;size:30" json:"billable_to_id"`
	CustomerID                          int       `gorm:"column:customer_id;size:30" json:"customer_id"`
	PayType                             string    `gorm:"column:pay_type;size:100" json:"pay_type"`
	FromEra                             string    `gorm:"column:from_era;size:1" json:"from_era"`
	QueueNumber                         string    `gorm:"column:queue_number;size:100" json:"queue_number"`
	ArrivalTime                         string    `gorm:"column:arrival_time;size:100" json:"arrival_time"`
	ServiceMileage                      string    `gorm:"column:service_mileage;size:100" json:"service_mileage"`
	LeaveCar                            string    `gorm:"column:leave_car;size:1" json:"leave_car"`
	Storing                             string    `gorm:"column:storing;size:1" json:"storing"`
	EraNumber                           string    `gorm:"column:era_number;size:100" json:"era_number"`
	EraExpiredDate                      time.Time `gorm:"column:era_expired_date;type:datetime" json:"era_expired_date"`
	Unregister                          string    `gorm:"column:unregister;size:1" json:"unregister"`
	ContactPersonName                   string    `gorm:"column:contact_person_name;size:100" json:"contact_person_name"`
	ContactPersonPhone                  string    `gorm:"column:contact_person_phone;size:100" json:"contact_person_phone"`
	ContactPersonMobile                 string    `gorm:"column:contact_person_mobile;size:100" json:"contact_person_mobile"`
	ContactPersonContactVia             string    `gorm:"column:contact_person_contact_via;size:100" json:"contact_person_contact_via"`
	ContractServiceSystemNumber         int       `gorm:"column:contract_service_system_number;size:30" json:"contract_service_system_number"`
	AgreementNumberGeneralRepairID      int       `gorm:"column:agreement_number_general_repair_id;size:30" json:"agreement_number_general_repair_id"`
	AgreementNumberBodyRepairID         int       `gorm:"column:agreement_number_body_repair_id;size:30" json:"agreement_number_body_repair_id"`
	BookingSystemNumber                 int       `gorm:"column:booking_system_number;size:30" json:"booking_system_number"`
	EstimationSystemNumber              int       `gorm:"column:estimation_system_number;size:30" json:"estimation_system_number"`
	PDISystemNumber                     float32   `gorm:"column:pdi_system_number" json:"pdi_system_number"`
	PDIDocumentNumber                   string    `gorm:"column:pdi_document_number;size:100" json:"pdi_document_number"`
	PDILineNumberID                     float32   `gorm:"column:pdi_line_number_id" json:"pdi_line_number_id"`
	ServiceRequestSystemNumber          int       `gorm:"column:service_request_system_number;size:30" json:"service_request_system_number"`
	CampaignID                          int       `gorm:"column:campaign_id;size:30" json:"campaign_id"`
	InsurancePolicyNumber               string    `gorm:"column:insurance_policy_number;size:100" json:"insurance_policy_number"`
	InsuranceExpiredDate                time.Time `gorm:"column:insurance_expired_date;type:datetime" json:"insurance_expired_date"`
	InsuranceClaimNumber                string    `gorm:"column:insurance_claim_number;size:100" json:"insurance_claim_number"`
	InsurancePersonInCharge             string    `gorm:"column:insurance_person_in_charge;size:100" json:"insurance_person_in_charge"`
	InsuranceOwnRisk                    string    `gorm:"column:insurance_own_risk;size:100" json:"insurance_own_risk"`
	InsuranceWorkOrderNumber            string    `gorm:"column:insurance_work_order_number;size:100" json:"insurance_work_order_number"`
	TotalPackage                        float32   `gorm:"column:total_package" json:"total_package"`
	TotalOperation                      float32   `gorm:"column:total_operation" json:"total_operation"`
	TotalPart                           float32   `gorm:"column:total_part" json:"total_part"`
	TotalOil                            float32   `gorm:"column:total_oil" json:"total_oil"`
	TotalMaterial                       float32   `gorm:"column:total_material" json:"total_material"`
	TotalConsumableMaterial             float32   `gorm:"column:total_consumable_material" json:"total_consumable_material"`
	TotalSublet                         float32   `gorm:"column:total_sublet" json:"total_sublet"`
	TotalPriceAccessories               float32   `gorm:"column:total_price_accessories" json:"total_price_accessories"`
	TotalDiscount                       float32   `gorm:"column:total_discount" json:"total_discount"`
	TotalVAT                            float32   `gorm:"column:total_vat" json:"total_vat"`
	TotalAfterVAT                       float32   `gorm:"column:total_after_vat" json:"total_after_vat"`
	TotalPPH                            float32   `gorm:"column:total_pph" json:"total_pph"`
	DiscountRequestPercent              float32   `gorm:"column:discount_request_percent" json:"discount_request_percent"`
	DiscountRequestAmount               float32   `gorm:"column:discount_request_amount" json:"discount_request_amount"`
	VatTaxID                            int       `gorm:"column:vat_tax_id;size:30" json:"vat_tax_id"`
	VatTaxRate                          float32   `gorm:"column:vat_tax_rate" json:"vat_tax_rate"`
	DiscountStatusID                    int       `gorm:"column:discount_status_id;size:30" json:"discount_status_id"`
	LastApprovalByID                    int       `gorm:"column:last_approval_by_id;size:30" json:"last_approval_by_id"`
	LastApprovalDate                    time.Time `gorm:"column:last_approval_date;type:datetime" json:"last_approval_date"`
	Remark                              string    `gorm:"column:remark;size:100" json:"remark"`
	ForemanID                           int       `gorm:"column:foreman_id;size:30" json:"foreman_id"`
	ProductionHeadID                    int       `gorm:"column:production_head_id;size:30" json:"production_head_id"`
	EstimateTime                        float32   `gorm:"column:estimate_time" json:"estimate_time"`
	Notes                               string    `gorm:"column:notes;size:100" json:"notes"`
	Suggestion                          string    `gorm:"column:suggestion;size:100" json:"suggestion"`
	FSCouponNumber                      string    `gorm:"column:fs_coupon_number;size:100" json:"fs_coupon_number"`
	ServiceAdvisorID                    int       `gorm:"column:service_advisor_id;size:30" json:"service_advisor_id"`
	IncentiveDate                       time.Time `gorm:"column:incentive_date;type:datetime" json:"incentive_date"`
	WorkOrderCancelReason               string    `gorm:"column:work_order_cancel_reason;size:100" json:"work_order_cancel_reason"`
	InvoiceSystemNumber                 int       `gorm:"column:invoice_system_number;size:30" json:"invoice_system_number"`
	CurrencyID                          int       `gorm:"column:currency_id;size:30" json:"currency_id"`
	ATPMWarrantyClaimFormDocumentNumber string    `gorm:"column:atpm_warranty_claim_form_document_number;size:100" json:"atpm_warranty_claim_form_document_number"`
	ATPMWarrantyClaimFormDate           time.Time `gorm:"column:atpm_warranty_claim_form_date;type:datetime" json:"atpm_warranty_claim_form_date"`
	ATPMFreeServiceDocumentNumber       string    `gorm:"column:atpm_free_service_document_number;size:100" json:"atpm_free_service_document_number"`
	ATPMFreeServiceDate                 time.Time `gorm:"column:atpm_free_service_date;type:datetime" json:"atpm_free_service_date"`
	TotalAfterDiscount                  float32   `gorm:"column:total_after_discount" json:"total_after_discount"`
	ApprovalRequestNumberID             int       `gorm:"column:approval_request_number_id;size:30" json:"approval_request_number_id"`
	JournalSystemNumber                 int       `gorm:"column:journal_system_number;size:30" json:"journal_system_number"`
	ApprovalGatepassRequestNumber       int       `gorm:"column:approval_gatepass_request_number;size:30" json:"approval_gatepass_request_number"`
}

func (*WorkOrderHistory) TableName() string {
	return CreateWorkOrderHistoryTable
}
