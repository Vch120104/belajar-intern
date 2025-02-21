package transactionworkshoppayloads

import "time"

type BookingEstimationRequest struct {
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
	IsUnregistered                 bool      `json:"is_unregistered"`
	ContactPersonName              string    `json:"contact_person_name"`
	ContactPersonPhone             string    `json:"contact_person_phone"`
	ContactPersonMobile            string    `json:"contact_person_mobile"`
	ContactPersonViaId             int       `json:"contact_person_via_id"`
	InsurancePolicyNo              string    `json:"insurance_policy_no"`
	InsuranceExpiredDate           time.Time `json:"insurance_expired_date"`
	InsuranceClaimNo               string    `json:"insurance_claim_no"`
	InsurancePic                   string    `json:"insurance_pic"`
}

type BookingEstimationDetailRequest struct {
	EstimationDetailId                  int     `gorm:"column:estimation_detail_id" json:"estimation_detail_id" parent_entity:"trx_booking_estimation_detail" main_table:"trx_booking_estimation_detail"`
	EstimationSystemNumber              int     `gorm:"column:estimation_system_number" json:"estimation_system_number" parent_entity:"trx_booking_estimation_detail"`
	LineTypeId                          int     `gorm:"column:line_type_id" json:"line_type_id" parent_entity:"trx_booking_estimation_detail"`
	TransactionTypeId                   int     `gorm:"column:transaction_type_id" json:"transaction_type_id" parent_entity:"trx_booking_estimation_detail" `
	JobTypeId                           int     `gorm:"column:job_type_id" json:"job_type_id" parent_entity:"trx_booking_estimation_detail"`
	FrtQuantity                         float64 `gorm:"column:frt_quantity" json:"frt_quantity" parent_entity:"trx_booking_estimation_detail"`
	OperationItemId                     int     `gorm:"column:operation_item_id" json:"operation_item_id" parent_entity:"trx_booking_estimation_detail"`
	OperationItemCode                   string  `gorm:"column:operation_item_code" json:"operation_item_code" parent_entity:"trx_booking_estimation_detail"`
	OperationItemPrice                  float64 `gorm:"column:operation_item_price" json:"operation_item_price" parent_entity:"trx_booking_estimation_detail"`
	OperationItemDiscountAmount         float64 `gorm:"column:operation_item_discount_amount" json:"operation_item_discount_amount" parent_entity:"trx_booking_estimation_detail"`
	OperationItemDiscountPercent        float64 `gorm:"column:operation_item_discount_percent" json:"operation_item_discount_percent" parent_entity:"trx_booking_estimation_detail"`
	OperationItemDiscountRequestPercent float64 `gorm:"column:operation_item_discount_request_percent" json:"operation_item_discount_request_percent" parent_entity:"trx_booking_estimation_detail"`
	OperationItemDiscountRequestAmount  float64 `gorm:"column:operation_item_discount_request_amount" json:"operation_item_discount_request_amount" parent_entity:"trx_booking_estimation_detail"`
}
type BookingEstimationDetailResponse struct {
	EstimationSystemNumber             int     `json:"estimation_system_number"`
	EstimationDetailId                 int     `json:"estimation_detail_id"`
	EstimationDocumentNumber           string  `json:"estimation_document_number"`
	EstimationLine                     int     `json:"estimation_line"`
	LineTypeId                         int     `json:"line_type_id"`
	LineTypeCode                       string  `json:"line_type_code"`
	LineTypeName                       string  `json:"line_type_name"`
	TransactionTypeId                  int     `json:"transaction_type_id"`
	TransactionTypeCode                string  `json:"transaction_type_code"`
	JobTypeId                          int     `json:"job_type_id"`
	JobTypeCode                        string  `json:"job_type_code"`
	OperationItemId                    int     `json:"operation_item_id"`
	OperationItemCode                  string  `json:"operation_item_code"`
	Description                        string  `json:"description"`
	Uom                                string  `json:"uom"`
	FrtQuantity                        float64 `json:"frt_quantity"`
	OperationItemPrice                 float64 `json:"operation_item_price"`
	OperationItemDiscountAmount        float64 `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount float64 `json:"operation_item_discount_request_amount"`
}

type BookEstimRemarkRequest struct {
	BookingServiceRequest string `json:"booking_service_request"`
}

type BookEstimRemarkResponse struct {
	BookingServiceRequest string `json:"booking_service_request"`
	BookingDocumentNumber string `json:"booking_document_number"`
	BookingLine           int    `json:"booking_line"`
	IsActive              bool   `json:"is_active"`
}

type GetBookingById struct {
	BookingSystemNumber     int       `json:"booking_system_number"`
	Type                    string    `json:"type"`
	BatchNumber             string    `json:"batch_number"`
	CreatedBy               string    `json:"created_by"`
	BrandId                 int       `json:"brand_id"`
	ModelId                 int       `json:"model_id"`
	VehicleId               int       `json:"vehicle_id"`
	NoPolisi                string    `json:"no_polisi"`
	BillCustomerId          int       `json:"bill_customer_id"`
	IsUnregistered          bool      `json:"is_unregistered"`
	BookingNo               string    `json:"booking_no"`
	EstimationNo            string    `json:"estimation_no"`
	ServiceRequestNumber    string    `json:"service_request_number"`
	ServiceRequestCompanyId int       `json:"service_request_company_id"`
	PdiNo                   string    `json:"pdi_no"`
	ProfitCenter            int       `json:"profit_center"`
	DealerRepCodeId         int       `json:"dealer_respresentative_code_id"`
	Contactname             string    `json:"contact_name"`
	ContactPhoneNumber      string    `json:"contact_phone_number"`
	ContactMobilePhone      string    `json:"contact_mobile_phone"`
	ContactById             int       `json:"contact_by_id"`
	IsInsurance             bool      `json:"is_insurance"`
	InsurancePolicyNo       string    `json:"insurance_policy_no"`
	InsuranceExpireDate     time.Time `json:"insurance_expire_date"`
	ClaimNo                 string    `json:"claim_no"`
	InsurancePicName        string    `json:"insurance_pic_name"`
	CampaignId              int       `json:"campaign_id"`
	ContractSystemNumber    string    `json:"contact_system_number"`
	ContractExpireDate      time.Time `json:"contract_expire_date"`
	ContractDealer          string    `json:"contract_dealer_number"`
	AgreementNo             string    `json:"agreement_no"`
	AgreementExpireDate     time.Time `json:"agreement_expire_date"`
	AgreementDealer         string    `json:"agreement_dealer"`
	CallActivityId          int       `json:"call_activity_id"`
	EngineNo                string    `json:"engine_no"`
	VariantId               int       `json:"variant_id"`
	OptionId                int       `json:"option_id"`
	ColourId                int       `json:"colour_id"`
	DeliveryDate            time.Time `json:"delivery_date"`
	LastServiceDate         time.Time `json:"last_service_date"`
	LastServiceMileage      float64   `json:"last_service_mileage"`
	BillCustomerName        string    `json:"bill_customer_name"`
	BillAddressId           int       `json:"bill_address_id"`
	BillPhoneNumber         string    `json:"bill_pone_number"`
	BillfaxNo               string    `json:"bill_fax_no"`
	StnkName                string    `json:"stnk_name"`
	StnkAddressId           int       `json:"stnk_address_id"`
	BookingStatus           int       `json:"booking_status"`
	BookingDate             time.Time `json:"booking_date"`
	Status                  string    `json:"status"`
	Estimationdate          time.Time `json:"estimation_date"`
	DiscoutnAprrovalStatus  int       `json:"discount_approval_status"`
	ExpireDate              time.Time `json:"expire_date"`
	PdiRefNo                string    `json:"pdi_ref_no"`
	Stall                   string    `json:"stall"`
	TotalFrt                float64   ` json:"total_frt"`
	EstimatedTime           float64   `json:"estimated_time"`
	ServiceDate             time.Time `json:"service_date"`
	ServiceTime             float64   `json:"service_time"`
	TotalPackage            int       `json:"total_package"`
	TotalOperation          int       `json:"total_operation"`
	TotalPart               int       `json:"total_part"`
	TotalOil                int       `json:"total_oil"`
	TotalMaterial           int       `json:"total_material"`
	TotalSublet             int       `json:"total_sublet"`
	TotalAccessories        int       `json:"total_accessories"`
	SubTotal                float64   `json:"subtotal"`
	Discount                float64   `json:"discount"`
	Total                   float64   `json:"total"`
	VAT                     float64   `json:"vat"`
	GrandTotal              float64   `json:"grand_total"`
}

type GetBillAddress struct {
	BillAddressId  int    `json:"bill_address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}

type GetStnkAddress struct {
	StnkAddressId  int    `json:"stnk_address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}

type GetAllServiceBooking struct {
	ServiceRequestId int       `json:"service_request_id"`
	ServiceRequestNo string    `json:"service_request_no"`
	ServiceReqDate   time.Time `json:"service_req_date"`
	ReqCompanyName   string    `jaon:"req_company_name"`
	BrandId          int       `json:"brand_id"`
	ModelId          int       `json:"model_id"`
	VehicleId        int       `json:"vehicle_id"`
}

type GetAllDetailBooking struct {
	DetailId  int `json:"detail_id"`
	TypeId    int `json:"type_id"`
	TrxTypeId int `json:"transaction_type_id"`
	Number    int `json:"number"`
	ItemId    int `json:"item_id"`
}

type VehicleRemarkRequest struct {
	VehicleRequest string `json:"vehicle_request"`
}

type ReminderServicePost struct {
	BookingServiceReminder string `json:"booking_service_reminder"`
}

type BookEstimDetailReq struct {
	EstimationLineID               int        `json:"estimation_line_id"`
	EstimationLineCode             int        `json:"estimation_line_code"`
	BillID                         int        `json:"bill_id"`
	EstimationLineDiscountApproval int        `json:"estimation_line_discount_approval_status"`
	ItemOperationID                int        `json:"item_operation_id"`
	LineTypeID                     int        `json:"line_type_id"`
	PackageID                      int        `json:"package_id"`
	JobTypeID                      int        `json:"job_type_id"`
	FieldActionSystemNumber        int        `json:"field_action_system_number"`
	ApprovalRequestNumber          int        `json:"approval_request_number"`
	UOMID                          int        `json:"uom_id"`
	RequestDescription             string     `json:"request_description"`
	FRTQuantity                    float64    `json:"frt_quantity"`
	DiscountItemAmount             float64    `json:"discount_item_amount"`
	DiscountItemPercent            float64    `json:"discount_item_percent"`
	DiscountRequestPercent         float64    `json:"discount_request_percent"`
	DiscountRequestAmount          float64    `json:"discount_request_amount"`
	DiscountApprovalBy             string     `json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `json:"discount_approval_date"`
}

type BookEstimDetailUpdate struct {
	FRTQuantity            float64 `json:"frt_quantity"`
	DiscountRequestPercent float32 `json:"discount_request_percent"`
}

type BookEstimDetailPayloads struct {
	EstimationLineID               int        `json:"estimation_line_id"`
	EstimationLineCode             int        `json:"estimation_line_code"`
	EstimationSystemNumber         int        `json:"estimation_system_number"`
	BillID                         int        `json:"bill_id"`
	EstimationLineDiscountApproval int        `json:"estimation_line_discount_approval_status"`
	ItemOperationID                int        `json:"item_operation_id"`
	LineTypeID                     int        `json:"line_type_id"`
	PackageID                      int        `json:"package_id"`
	JobTypeID                      int        `json:"job_type_id"`
	FieldActionSystemNumber        int        `json:"field_action_system_number"`
	ApprovalRequestNumber          int        `json:"approval_request_number"`
	UOMID                          int        `json:"uom_id"`
	RequestDescription             string     `json:"request_description"`
	FRTQuantity                    float64    `json:"frt_quantity"`
	OperationItemPrice             float64    `json:"operation_item_price"`
	DiscountItemOperationAmount    float64    `json:"discount_item_operation_amount"`
	DiscountItemOperationPercent   float64    `json:"discount_item_operation_percent"`
	DiscountRequestPercent         float64    `json:"discount_request_percent"`
	DiscountRequestAmount          float64    `json:"discount_request_amount"`
	DiscountApprovalBy             string     `json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `json:"discount_approval_date"`
}

type BookEstimationPayloadsDiscount struct {
	PackageDiscount int `json:"package_discount"`
	Operation       int `json:"operation"`
	Sparepart       int `json:"sparepart"`
	Oil             int `json:"oil"`
	Material        int `json:"material"`
	Fee             int `json:"fee"`
	Accessories     int `json:"accessories"`
	Souvenir        int `json:"souvenir"`
}

type BookingEstimationPostPayloads struct {
	BrandId             int       `json:"brand_id"`
	ModelId             int       `json:"model_id"`
	VehicleId           int       `json:"vehicle_id"`
	DealerRepCodeId     int       `json:"dealer_rep_code_id"`
	ContactPersonName   string    `json:"contact_person_name"`
	ContactPersonPhone  string    `json:"contact_person_phne"`
	ContactPersonMobile string    `json:"contact_person_mobile"`
	ContactViaId        int       `json:"contact_person_via"`
	InsPolicyNo         string    `json:"insurance_policy_number"`
	InsExpireDate       time.Time `json:"insurance_expire_date"`
	InsClaimNo          string    `json:"insurance_claim_no"`
	InsPIC              string    `json:"insurace_pic"`
	CampagnCodeId       int       `json:"campaign_code_id"`
}

type BookingEstimationCalculationPayloads struct {
	EstimationSystemNumber           int       `json:"estimation_system_number"`
	BatchSystemNumber                int       `json:"batch_system_number"`
	DocumentStatusId                 int       `json:"document_status_id"`
	EstimationDiscountApprovalStatus int       `json:"estimation_discount_approval_status"`
	CompanyId                        int       `json:"company_id"`
	ApprovalRequestNumber            int       `json:"approval_request_number"`
	EstimationDoumentNumber          string    `json:"estimation_document_number"`
	EstimationDate                   time.Time `json:"estimation_date"`
	TotalPricePackage                float64   `json:"total_price_package"`
	TotalPriceOperation              float64   `json:"total_price_operation"`
	TotalPricePart                   float64   `json:"total_price_part"`
	TotalPriceOil                    float64   `json:"total_price_oil"`
	TotalPriceMaterial               float64   `json:"total_price_material"`
	TotalPriceConsumableMaterial     float64   `json:"total_price_consumable_material"`
	TotalSublet                      float64   `json:"total_sublet"`
	TotalPriceAccessories            float64   `json:"total_price_accessories"`
	TotalDiscount                    float64   `json:"total_discount"`
	TotalVat                         float64   `json:"total_vat"`
	TotalAfterVat                    float64   `json:"total_after_vat"`
	AdditionalDiscountRequestPercent float64   `json:"additional_discount_request_percent"`
	AdditionalDiscountRequestAmount  float64   `json:"additional_discount_request_amount"`
	VatTaxRate                       float64   `json:"vat_tax_rate"`
	DiscountApprovalBy               string    `json:"discount_approval_by"`
	DiscountApprovalDate             time.Time `json:"discount_approval_date"`
	TotalAfterDiscount               float64   `json:"total_after_discount"`
}

type BookEstimDetailPayloadsOperation struct {
	LineTypeid        int     `json:"line_type_id"`
	TransactionTypeId int     `json:"transaction_type_id"`
	OperationId       int     `json:"operation_id"`
	OperationName     string  `json:"operation_name"`
	Quantity          int     `json:"quantity"`
	Price             float64 `json:"price"`
	SubTotal          float64 `json:"subtotal"`
	OriginalDiscount  float64 `json:"original_discount"`
	ProposalDiscount  float64 `json:"proposal_discount"`
	Total             float64 `json:"total"`
}

type TransactionTypePayloads struct {
	TransactionTypeId   int    `json:"transaction_type_id"`
	TransactionTypeCode string `jon:"transaction_type_code"`
}

type BookEstimDetailPayloadsItem struct {
	LineTypeid        int     `json:"line_type_id"`
	TransactionTypeId int     `json:"transaction_type_id"`
	ItemId            int     `json:"item_id"`
	ItemName          string  `json:"item_name"`
	Quantity          int     `json:"quantity"`
	Price             float64 `json:"price"`
	SubTotal          float64 `json:"subtotal"`
	OriginalDiscount  float64 `json:"original_discount"`
	ProposalDiscount  float64 `json:"proposal_discount"`
	Total             float64 `json:"total"`
}

type VehicleDetailPayloads struct {
	EngineNo           string `json:"engine_no"`
	VariantId          int    `json:"variant_id"`
	OptionId           int    `json:"option_id"`
	ColourId           int    `json:"colour_id"`
	DeliveryDate       string `json:"delivery_date"`
	LastServiceDate    string `json:"last_service_date"`
	LastServiceMileage string `json:"last_service_mileage"`
	StnkName           string `json:"stnk_name"`
	StnkAddress        string `json:"stnk_address"`
	BillingName        string `json:"billing_name"`
	BillingAddress     string `json:"billing_address"`
	BillingZipNo       string `json:"billing_zip_number"`
}

type BillingDetail struct {
	Name       string `json:"name"`
	AddressId  int    `json:"address_id"`
	VillageId  int    `json:"villgae_id"`
	DistrictId int    `json:"district_id"`
	CityId     int    `json:"cistr_id"`
	ZipCode    string `json:"zip_code"`
	PhoneNo    string `json:"phone_number"`
	FaxNumber  string `json:"fax_number"`
}

type StnkDetail struct {
	STNKName    string `json:"stnk_name"`
	STNKAddress string `json:"stnk_address"`
	Village     string `json:"village"`
	Disrict     string `json:"dictrict"`
	City        string `json:"city"`
	Province    string `json:"province"`
	ZipCode     string `json:"zip_code"`
}

type VehicleTnkb struct {
	VehicleId        int    `json:"vehicle_id"`
	VehicleBrandId   int    `json:"vehicle_brand_id"`
	VehicleModelId   int    `json:"vehicle_model_id"`
	VehicleVariantId int    `json:"vehicle_variant_id"`
	Tnkb             string `json:"vehicle_registration_certificate_tnkb"`
}

type GetAllBookEstim struct {
	BatchSystemNo       int       `json:"batch_system_number"`
	BookingSystemNumber int       `json:"booking_system_number"`
	ServiceDate         time.Time `json:"service_date"`
	ServiceTime         float64   `json:"service_time"`
	BookingStatusID     int       `json:"booking_status_id"`
	EstimationNo        string    `json:"estimation_number"`
	EstimationDate      time.Time `json:"estimation_date"`
	EstimationStatus    int       `json:"estimation_status"`
}

type BookEstimationAllocation struct {
	BookingStatusID       int        `json:"booking_status_id"`
	CompanyID             int        `json:"company_id"`
	PdiSystemNumber       int        `json:"pdi_system_number"`
	BookingDocumentNumber string     `json:"booking_document_number"`
	BookingDate           *time.Time `json:"booking_date"`
	BookingStall          string     `json:"booking_stall"`
	BookingReminderDate   *time.Time `json:"booking_reminder_date"`
	BookingServiceDate    *time.Time `json:"booking_service_date"`
	BookingServiceTime    float64    `json:"booking_service_time"`
	BookingEstimationTime float64    `json:"booking_estimation_time"`
}

type BookingEstimationFirstContractService struct {
	ContractServiceSystemNumber int    `json:"contract_service_system_number"`
	EstimationDiscountStatus    string `json:"estimation_disount_status"`
	BookingSystemNumber         int    `json:"booking_system_number"`
	EstimationDocumentNumber    string `json:"estimation_document_number"`
	BrandId                     string `json:"brand_id"`
	ProfitCenterId              string `json:"profit_center_id"`
	ModelId                     string `json:"model_id"`
	CompanyId                   int    `json:"company_id"`
	EstimationSystemNumber      int    `json:"estimation_system_number"`
}

type ContractService struct {
	ItemOperationId     int     `json:"item_operation_id"`
	LineTypeId          int     `json:"line_type_id"`
	Description         string  `json:"description"`
	FrtQuantity         int     `json:"frt_quantity"`
	ItemPrice           float64 `json:"item_price"`
	ItemDiscountPercent float64 `json:"item_discount_percent"`
}

type PackageForDetail struct {
	ItemOperationId     int     `json:"item_operation_id"`
	LineTypeId          int     `json:"line_type_id"`
	ItemOrOperationName string  `json:"item_or_operation_name"` // This will be populated based on the CASE statement
	FrtQuantity         float64 `json:"frt_qty"`
	CurrencyId          int     `json:"currency_id"`
	JobTypeId           int     `json:"job_type_id"`
	TransactionTypeId   int     `json:"transacion_type_id"`
	BillId              int     `json:"bill_id"`
}

type CompanyReference struct {
	CurrencyId                int     `json:"currency_id"`
	CoaGroupId                int     `json:"coa_group_id"`
	OperationDiscountOuterKpp float64 `json:"operation_discount_outer_kpp"`
	MarginOuterKpp            float64 `json:"margin_outer_kpp"`
	AdjustmentReasonId        int     `json:"adjustment_reason_id"`
	LeadTimeUnitEtd           int     `json:"lead_time_unit_etd"`
	BankAccReceiveCompanyId   int     `json:"bank_acc_receive_company_id"`
	UnitWarehouseId           int     `json:"unit_warehouse_id"`
	TimeDifference            float64 `json:"time_difference"`
	UseDms                    bool    `json:"use_dms"`
	UseJpcb                   bool    `json:"use_jpcb"`
	CheckMonthEnd             bool    `json:"check_month_end"`
	IsDistributor             bool    `json:"is_distributor"`
	WithVat                   bool    `json:"with_vat"`
}

type WorkorderTransactionType struct {
	WorkOrderTransactionTypeName string `json:"work_order_transaction_type_name"`
	WorkOrderTransactionTypeId   int    `json:"work_order_transaction_type_id"`
	IsActive                     bool   `json:"is_active"`
	WorkOrderTransactionTypeCode string `json:"work_order_transaction_type_code"`
}

type DocumentStatus struct {
	IsActive                  bool   `json:"is_active"`
	DocumentStatusId          int    `json:"document_status_id"`
	DocumentStatusCode        string `json:"document_status_code"`
	DocumentStatusDescription string `json:"document_status_description"`
}

type ApprovalStatus struct {
	ApprovalStatusId   int    `json:"approval_status_id"`
	ApprovalStatusCode string `json:"approval_status_code"`
}

type PdiServiceRequest struct {
	ContactPersonName   string `json:"contact_person_name"`
	ContactPersonPhone  string `json:"contact_person_phone"`
	ContactPersonMobile string `json:"contact_person_mobile"`
	ContactPersonViaId  int    `json:"contact_person_via_id"`
}

type BookingEstimationNormalRequest struct {
	BrandId                int  `json:"brand_id"`
	ModelId                int  `json:"model_id"`
	VariantId              int  `json:"variant_id"`
	VehicleId              int  `json:"vehicle_id"`
	Unregistered           bool `json:"unregistered"`
	ProfitCenterId         int  `json:"profit_center_id"`
	DealerRepresentativeId int  `json:"dealer_representative_id"`
	OrderId                int  `json:"order_id"`
	CompanyId              int  `json:"company_id"`
	CustomerId             int  `json:"customer_id"`

	// Customer contact information
	NameCust   string `json:"name_customer"`
	PhoneCust  string `json:"phone_customer"`
	MobileCust string `json:"mobile_customer"`
	ContactVia string `json:"contact_via"`

	// Insurance information
	InsuranceCheck    bool      `json:"insurance_check"`
	InsurancePolicyNo string    `json:"insurance_policy_no"`
	InsuranceExpired  time.Time `json:"insurance_expired_date"`
	InsuranceClaimNo  string    `json:"insurance_claim_no"`
	InsurancePic      string    `json:"insurance_pic"`

	// Agreement information
	AgreementId      int       `json:"agreement_id"`
	AgreementNumber  string    `json:"agreement_number"`
	AgreementExpired time.Time `json:"agreement_expired_date"`
	AgreementDealer  string    `json:"agreement_dealer"`

	// Campaign information
	CampaignId int `json:"campaign_id"`

	// Contract information
	ContractSystemNumber int       `json:"contract_system_number"`
	ContractExpiredDate  time.Time `json:"contract_expired_date"`
	ContractDealer       string    `json:"contract_dealer"`

	CallActivityId int `json:"call_activity_id"`

	// page 2
	Stall         string  `json:"stall"`
	EstimatedTime float64 `json:"estimated_time"`
}

type BookingEstimationSaveRequest struct {
	DealerRepresentativeId int `json:"dealer_representative_id"`
	OrderId                int `json:"order_id"`
	CompanyId              int `json:"company_id"`

	// Customer contact information
	NameCust     string `json:"name_customer"`
	PhoneCust    string `json:"phone_customer"`
	MobileCust   string `json:"mobile_customer"`
	ContactViaId int    `json:"contact_person_via"`

	// Insurance information
	InsuranceCheck    bool      `json:"insurance_check"`
	InsurancePolicyNo string    `json:"insurance_policy_no"`
	InsuranceExpired  time.Time `json:"insurance_expired_date"`
	InsuranceClaimNo  string    `json:"insurance_claim_no"`
	InsurancePic      string    `json:"insurance_pic"`

	// Campaign information
	CampaignId int `json:"campaign_id"`

	// page 2
	Stall         string  `json:"stall"`
	EstimatedTime float64 `json:"estimated_time"`
}

type BookingEstimationResponseDetail struct {
	BatchSystemNumber          int       `json:"batch_system_number"`
	BookingEstimationBatchDate time.Time `json:"booking_estimation_batch_date"`

	BrandId                int    `json:"brand_id"`
	BrandName              string `json:"brand_name"`
	ModelId                int    `json:"model_id"`
	ModelDescription       string `json:"model_description"`
	VariantId              int    `json:"variant_id"`
	VariantDescription     string `json:"variant_description"`
	VehicleId              int    `json:"vehicle_id"`
	Unregistered           bool   `json:"unregistered"`
	ProfitCenterId         int    `json:"profit_center_id"`
	DealerRepresentativeId int    `json:"dealer_representative_id"`
	OrderId                int    `json:"order_id"`
	CompanyId              int    `json:"company_id"`
	CustomerId             int    `json:"customer_id"`

	// Customer contact information
	NameCust   string `json:"name_customer"`
	PhoneCust  string `json:"phone_customer"`
	MobileCust string `json:"mobile_customer"`
	ContactVia string `json:"contact_via"`

	// Insurance information
	InsuranceCheck    bool      `json:"insurance_check"`
	InsurancePolicyNo string    `json:"insurance_policy_no"`
	InsuranceExpired  time.Time `json:"insurance_expired_date"`
	InsuranceClaimNo  string    `json:"insurance_claim_no"`
	InsurancePic      string    `json:"insurance_pic"`

	// Agreement information
	AgreementId      int       `json:"agreement_id"`
	AgreementNumber  string    `json:"agreement_number"`
	AgreementExpired time.Time `json:"agreement_expired_date"`
	AgreementDealer  string    `json:"agreement_dealer"`

	// Campaign information
	CampaignId int `json:"campaign_id"`

	// Contract information
	ContractSystemNumber int       `json:"contract_system_number"`
	ContractExpiredDate  time.Time `json:"contract_expired_date"`
	ContractDealer       string    `json:"contract_dealer"`

	CallActivityId int `json:"call_activity_id"`

	// page 2
	VehicleCode string `json:"vehicle_code"`
	VehicleTnkb string `json:"vehicle_tnkb"`

	BookingStatusId       int       `json:"booking_status_id"`
	BookingDocumentNumber string    `json:"booking_document_number"`
	BookingDate           time.Time `json:"booking_date"`

	EstimationStatusId       int       `json:"estimation_status_id"`
	EstimationDocumentNumber string    `json:"estimation_document_number"`
	EstimationDate           time.Time `json:"estimation_date"`
	DiscountApprovalStatusId int       `json:"discount_approval_status_id"`
	ExpireDate               time.Time `json:"expire_date"`
	PdiRefNo                 string    `json:"pdi_ref_no"`
	Stall                    string    `json:"stall"`

	TotalFrt      float64   `json:"total_frt"`
	EstimatedTime float64   `json:"estimated_time"`
	ServiceDate   time.Time `json:"service_date"`
	ServiceTime   float64   `json:"service_time"`

	BookingEstimationCampaign  BookingEstimationCampaignDetail `json:"booking_estimation_campaign"`
	BookingEstimationContract  BookingEstimationContractDetail `json:"booking_estimation_contract"`
	BookingEstimationAgreement BookingEstimationAgreement      `json:"booking_estimation_agreement"`
}

type BookingEstimationCampaignDetail struct {
	DataCampaign []BookingEstimationCampaignResponse `json:"data"`
}
type BookingEstimationCampaignResponse struct {
	CampaignId         int       `json:"campaign_id"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	CampaignPeriodFrom time.Time `json:"campign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
}

type BookingEstimationContractDetail struct {
	DataContract []BookingEstimationContractsResponse `json:"data"`
}

type BookingEstimationContractsResponse struct {
	ContractServiceSystemNumber   int       `json:"contract_service_system_number"`
	ContractServiceDocumentNumber int       `json:"contract_service_document_number"`
	ContractServiceExpiredDate    time.Time `json:"contract_service_expired_date"`
	CompanyId                     int       `json:"company_id"`
}

type BookingEstimationAgreement struct {
	DataAgreement []BookingEstimationAgreementResponse `json:"data"`
}

type BookingEstimationAgreementResponse struct {
	AgreementId       int       `json:"agreement_id"`
	AgreementCode     string    `json:"agreement_code"`
	AgreementDateFrom time.Time `json:"agreement_date_from"`
	AgreementDateTo   time.Time `json:"agreement_date_to"`
	CompanyId         int       `json:"company_id"`
}
