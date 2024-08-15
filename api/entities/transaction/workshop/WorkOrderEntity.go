package transactionworkshopentities

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"time"
)

const TableNameWorkOrder = "trx_work_order"

type WorkOrder struct {
	WorkOrderSystemNumber              int                                     `gorm:"column:work_order_system_number;size:30;primaryKey" json:"work_order_system_number"`
	CompanyId                          int                                     `gorm:"column:company_id;size:30;" json:"company_id"`
	WorkOrderDocumentNumber            string                                  `gorm:"column:work_order_document_number;size:50;" json:"work_order_document_number"`
	WorkOrderStatusId                  int                                     `gorm:"column:work_order_status_id;size:30;" json:"work_order_status_id"`
	WorkOrderDate                      *time.Time                              `gorm:"column:work_order_date;default:null" json:"work_order_date"`
	WorkOrderCloseDate                 *time.Time                              `gorm:"column:work_order_close_date;default:null" json:"work_order_close_date"`
	WorkOrderTypeId                    int                                     `gorm:"column:work_order_type_id;size:30;" json:"work_order_type_id"`
	WorkOrderRepeatedSystemNumber      int                                     `gorm:"column:work_order_repeated_system_number;size:30;" json:"work_order_repeated_system_number"`
	WorkOrderRepeatedDocumentNumber    string                                  `gorm:"column:work_order_repeated_document_number;size:50;" json:"work_order_repeated_document_number"`
	AffiliatedCompany                  int                                     `gorm:"column:afiliated_company;size:30;" json:"afiliated_company"`
	ProfitCenterId                     int                                     `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	BrandId                            int                                     `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId                            int                                     `gorm:"column:model_id;size:30;" json:"model_id"`
	ServiceSite                        string                                  `gorm:"column:service_site;size:50;" json:"service_site"`
	VariantId                          int                                     `gorm:"column:variant_id;size:30;" json:"variant_id"`
	VehicleChassisNumber               string                                  `gorm:"column:vehicle_chassis_number;size:50;" json:"vehicle_chassis_number"`
	VehicleId                          int                                     `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	BillableToId                       int                                     `gorm:"column:billable_to_id;size:30;" json:"billable_to_id"`
	CustomerId                         int                                     `gorm:"column:customer_id;size:30;" json:"customer_id"`
	PayType                            string                                  `gorm:"column:pay_type;size:50;" json:"pay_type"`
	FromEra                            bool                                    `gorm:"column:from_era;default:false;" json:"from_era"`
	QueueNumber                        int                                     `gorm:"column:queue_number;size:50;" json:"queue_number"`
	ArrivalTime                        *time.Time                              `gorm:"column:arrival_time;default:null" json:"arrival_time"`
	ServiceMileage                     int                                     `gorm:"column:service_mileage;size:50;" json:"service_mileage"`
	LeaveCar                           bool                                    `gorm:"column:leave_car;default:false;" json:"leave_car"`
	Storing                            bool                                    `gorm:"column:storing;default:false;" json:"storing"`
	EraNumber                          string                                  `gorm:"column:era_number;size:50;" json:"era_number"`
	EraExpiredDate                     *time.Time                              `gorm:"column:era_expired_date;default:null" json:"era_expired_date"`
	Unregister                         bool                                    `gorm:"column:unregister;size:50;" json:"unregister"`
	ContactPersonName                  string                                  `gorm:"column:contact_person_name;size:50;" json:"contact_person_name"`
	ContactPersonPhone                 string                                  `gorm:"column:contact_person_phone;size:50;" json:"contact_person_phone"`
	ContactPersonMobile                string                                  `gorm:"column:contact_person_mobile;size:50;" json:"contact_person_mobile"`
	ContactPersonMobileAlternative     string                                  `gorm:"column:contact_person_mobile_alternative;size:50;" json:"contact_person_mobile_alternative"`
	ContactPersonMobileDriver          string                                  `gorm:"column:contact_person_mobile_driver;size:50;" json:"contact_person_mobile_driver"`
	ContactPersonContactVia            string                                  `gorm:"column:contact_person_contact_via;size:50;" json:"contact_person_contact_via"`
	ContractServiceSystemNumber        int                                     `gorm:"column:contract_service_system_number;size:30;" json:"contract_service_system_number"`
	AgreementGeneralRepairId           int                                     `gorm:"column:agrement_general_repair_id;size:30;" json:"agreement_general_repair_id"`
	AgreementBodyRepairId              int                                     `gorm:"column:agreement_body_repair_id;size:30;" json:"agreement_body_repair_id"`
	BookingSystemNumber                int                                     `gorm:"column:booking_system_number;size:30;" json:"booking_system_number"`
	EstimationSystemNumber             int                                     `gorm:"column:estimation_system_number;size:30;" json:"estimation_system_number"`
	PDISystemNumber                    *float64                                `gorm:"column:pdi_system_number;default:null" json:"pdi_system_number"`
	PDIDocumentNumber                  string                                  `gorm:"column:pdi_document_number;size:50;" json:"pdi_document_number"`
	PDILineNumber                      int                                     `gorm:"column:pdi_line_number;size:30;" json:"pdi_line_number"`
	ServiceRequestSystemNumber         int                                     `gorm:"column:service_request_system_number;size:30;" json:"service_request_system_number"`
	CampaignId                         int                                     `gorm:"column:campaign_id;size:30;" json:"campaign_id"`
	CampaignCode                       int                                     `gorm:"column:campaign_code;size:30;" json:"campaign_code"`
	InsuranceCheck                     bool                                    `gorm:"column:insurance_check;default:false;" json:"insurance_check"`
	InsurancePolicyNumber              string                                  `gorm:"column:insurance_policy_number;size:50;" json:"insurance_policy_number"`
	InsuranceExpiredDate               *time.Time                              `gorm:"column:insurance_expired_date;default:null" json:"insurance_expired_date"`
	InsuranceClaimNumber               string                                  `gorm:"column:insurance_claim_number;size:50;" json:"insurance_claim_number"`
	InsurancePersonInCharge            string                                  `gorm:"column:insurance_person_in_charge;size:50;" json:"insurance_person_in_charge"`
	InsuranceOwnRisk                   *float64                                `gorm:"column:insurance_own_risk;default:null" json:"insurance_own_risk"`
	InsuranceWorkOrderNumber           string                                  `gorm:"column:insurance_work_order_number;size:50;" json:"insurance_work_order_number"`
	TotalPackage                       *float64                                `gorm:"column:total_package;default:null" json:"total_package"`
	TotalOperation                     *float64                                `gorm:"column:total_operation;default:null" json:"total_operation"`
	TotalPart                          *float64                                `gorm:"column:total_part;default:null" json:"total_part"`
	TotalOil                           *float64                                `gorm:"column:total_oil;default:null" json:"total_oil"`
	TotalMaterial                      *float64                                `gorm:"column:total_material;default:null" json:"total_material"`
	TotalConsumableMaterial            *float64                                `gorm:"column:total_consumable_material;default:null" json:"total_consumable_material"`
	TotalSublet                        *float64                                `gorm:"column:total_sublet;default:null" json:"total_sublet"`
	TotalPriceAccessories              *float64                                `gorm:"column:total_price_accessories;default:null" json:"total_price_accessories"`
	TotalDiscount                      *float64                                `gorm:"column:total_discount;default:null" json:"total_discount"`
	Total                              *float64                                `gorm:"column:total;default:null" json:"total"`
	TotalVAT                           *float64                                `gorm:"column:total_vat;default:null" json:"total_vat"`
	TotalAfterVAT                      *float64                                `gorm:"column:total_after_vat;default:null" json:"total_after_vat"`
	TotalPPH                           *float64                                `gorm:"column:total_pph;default:null" json:"total_pph"`
	DiscountRequestPercent             *float64                                `gorm:"column:discount_request_percent;default:null" json:"discount_request_percent"`
	DiscountRequestAmount              *float64                                `gorm:"column:discount_request_amount;default:null" json:"discount_request_amount"`
	TaxId                              int                                     `gorm:"column:tax_id;size:30;" json:"tax_id"`
	VATTaxRate                         *float64                                `gorm:"column:vat_tax_rate;default:null" json:"vat_tax_rate"`
	AdditionalDiscountStatusApprovalId int                                     `gorm:"column:additional_discount_status_approval;size:30;" json:"additional_discount_status_approval"`
	LastApprovalBy                     int                                     `gorm:"column:last_approval_by_id;size:30;" json:"last_approval_by_id"`
	LastApprovalDate                   *time.Time                              `gorm:"column:last_approval_date;default:null" json:"last_approval_date"`
	Remark                             string                                  `gorm:"column:remark;size:50;" json:"remark"`
	Foreman                            int                                     `gorm:"column:foreman_id;size:30;" json:"foreman_id"`
	ProductionHead                     int                                     `gorm:"column:production_head_id;size:30;" json:"production_head_id"`
	EstTime                            *float64                                `gorm:"column:estimate_time;default:null" json:"estimate_time"`
	Notes                              string                                  `gorm:"column:notes;size:50;" json:"notes"`
	Suggestion                         string                                  `gorm:"column:suggestion;size:50;" json:"suggestion"`
	FSCouponNo                         string                                  `gorm:"column:fs_coupon_number;size:50;" json:"fs_coupon_number"`
	ServiceAdvisor                     int                                     `gorm:"column:service_advisor_id;size:30;" json:"service_advisor_id"`
	IncentiveDate                      *time.Time                              `gorm:"column:incentive_date;default:null" json:"incentive_date"`
	WOCancelReason                     string                                  `gorm:"column:work_order_cancel_reason;size:50;" json:"work_order_cancel_reason"`
	InvoiceSystemNumber                int                                     `gorm:"column:invoice_system_number;size:30;" json:"invoice_system_number"`
	CurrencyId                         int                                     `gorm:"column:currency_id;size:30;" json:"currency_id"`
	ATPMWCFDocNo                       string                                  `gorm:"column:atpm_warranty_claim_form_document_number;size:50;" json:"atpm_warranty_claim_form_document_number"`
	ATPMWCFDate                        *time.Time                              `gorm:"column:atpm_warranty_claim_form_date;default:null" json:"atpm_warranty_claim_form_date"`
	ATPMFSDocNo                        string                                  `gorm:"column:atpm_free_service_document_number;size:50;" json:"atpm_free_service_document_number"`
	ATPMFSDate                         *time.Time                              `gorm:"column:atpm_free_service_date;default:null" json:"atpm_free_service_date"`
	TotalAfterDisc                     *float64                                `gorm:"column:total_after_discount;default:null" json:"total_after_discount"`
	ApprovalReqNo                      int                                     `gorm:"column:approval_request_number;size:30;" json:"approval_request_number"`
	JournalSysNo                       int                                     `gorm:"column:journal_system_number;size:30;" json:"journal_system_number"`
	ApprovalGatepassReqNo              int                                     `gorm:"column:approval_gatepass_request_number;size:30;" json:"approval_gatepass_request_number"`
	DPAmount                           *float64                                `gorm:"column:downpayment_amount;default:null" json:"downpayment_amount"`
	DPPayment                          *float64                                `gorm:"column:downpayment_payment;default:null" json:"downpayment_payment"`
	DPPaymentAllocated                 *float64                                `gorm:"column:downpayment_payment_allocated;default:null" json:"downpayment_payment_allocated"`
	DPPaymentVAT                       *float64                                `gorm:"column:downpayment_payment_vat;default:null" json:"downpayment_payment_vat"`
	DPAllocToInv                       *float64                                `gorm:"column:downpayment_payment_to_invoice;default:null" json:"downpayment_payment_to_invoice"`
	DPVATAllocToInv                    *float64                                `gorm:"column:downpayment_payment_vat_to_invoice;default:null" json:"downpayment_payment_vat_to_invoice"`
	JournalOverpaySysNo                int                                     `gorm:"column:journal_overpay_system_number;size:30;" json:"journal_overpay_system_number"`
	DPOverpay                          *float64                                `gorm:"column:downpayment_overpay;default:null" json:"downpayment_overpay"`
	SiteTypeId                         int                                     `gorm:"column:work_order_site_type_id;size:30;" json:"work_order_site_type_id"`
	CostCenterId                       int                                     `gorm:"column:cost_center_id;size:30;" json:"cost_center_id"`
	PromiseDate                        *time.Time                              `gorm:"column:promise_date;default:null" json:"promise_date"`
	PromiseTime                        *time.Time                              `gorm:"column:promise_time;default:null" json:"promise_time"`
	CarWash                            bool                                    `gorm:"column:car_wash;default:false;" json:"car_wash"`
	JobOnHoldReason                    string                                  `gorm:"column:job_on_hold_reason;size:50;" json:"job_on_hold_reason"`
	CustomerExpress                    bool                                    `gorm:"column:customer_express;default:false;" json:"customer_express"`
	CPTitlePrefix                      string                                  `gorm:"column:contact_person_title_prefix;size:50;" json:"contact_person_title_prefix"`
	WorkOrderDetail                    []WorkOrderDetail                       `gorm:"foreignKey:WorkOrderSystemNumber;references:WorkOrderSystemNumber" json:"work_order_detail"`
	SupplySlip                         transactionsparepartentities.SupplySlip `gorm:"foreignkey:WorkOrderSystemNumber;references:WorkOrderSystemNumber"`
}

func (*WorkOrder) TableName() string {
	return TableNameWorkOrder
}
