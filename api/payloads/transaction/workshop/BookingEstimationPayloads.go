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

type BookEstimRemarkRequest struct {
	BookingServiceRequest string `json:"booking_service_request"`
}

type GetBookingById struct {
	BookingId               int       `json:"booking_id"`
	Type                    string    `json:"type"`
	BatchNumber             string    `json:"batch_number"`
	CreatedBy               string    `json:"created_by"`
	BrandName               string    `json:"brand_name"`
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
	ContractServiceNo       string    `json:"contact_service_no"`
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
	EstimationSystemNumber         int        `json:"estimation_system_number"`
	BillID                         int        `json:"bill_id"`
	EstimationLineDiscountApproval int        `json:"estimation_line_discount_approval_status"`
	ItemID                         int        `json:"item_id"`
	OperationId                    int        `json:"operation_id"`
	LineTypeID                     int        `json:"line_type_id"`
	PackageID                      int        `json:"package_id"`
	JobTypeID                      int        `json:"job_type_id"`
	FieldActionSystemNumber        int        `json:"field_action_system_number"`
	ApprovalRequestNumber          int        `json:"approval_request_number"`
	UOMID                          int        `json:"uom_id"`
	RequestDescription             string     `json:"request_description"`
	FRTQuantity                    float64    `json:"frt_quantity"`
	OperationItemPrice             float64    `json:"operation_item_price"`
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

type BookEstimItemPayloads struct {
	EstimationLineID               int        `json:"estimation_line_id"`
	EstimationLineCode             int        `json:"estimation_line_code"`
	EstimationSystemNumber         int        `json:"estimation_system_number"`
	BillID                         int        `json:"bill_id"`
	EstimationLineDiscountApproval int        `json:"estimation_line_discount_approval_status"`
	ItemID                         int        `json:"item_id"`
	LineTypeID                     int        `json:"line_type_id"`
	PackageID                      int        `json:"package_id"`
	JobTypeID                      int        `json:"job_type_id"`
	FieldActionSystemNumber        int        `json:"field_action_system_number"`
	ApprovalRequestNumber          int        `json:"approval_request_number"`
	UOMID                          int        `json:"uom_id"`
	RequestDescription             string     `json:"request_description"`
	FRTQuantity                    float64    `json:"frt_quantity"`
	OperationItemPrice             float64    `json:"operation_item_price"`
	DiscountItemAmount             float64    `json:"discount_item_amount"`
	DiscountItemPercent            float64    `json:"discount_item_percent"`
	DiscountRequestPercent         float64    `json:"discount_request_percent"`
	DiscountRequestAmount          float64    `json:"discount_request_amount"`
	DiscountApprovalBy             string     `json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `json:"discount_approval_date"`
}

type BookEstimOperationPayloads struct {
	EstimationLineID               int        `json:"estimation_line_id"`
	EstimationLineCode             int        `json:"estimation_line_code"`
	EstimationSystemNumber         int        `json:"estimation_system_number"`
	BillID                         int        `json:"bill_id"`
	EstimationLineDiscountApproval int        `json:"estimation_line_discount_approval_status"`
	OperationId                    int        `json:"operation_id"`
	LineTypeID                     int        `json:"line_type_id"`
	PackageID                      int        `json:"package_id"`
	JobTypeID                      int        `json:"job_type_id"`
	FieldActionSystemNumber        int        `json:"field_action_system_number"`
	ApprovalRequestNumber          int        `json:"approval_request_number"`
	UOMID                          int        `json:"uom_id"`
	RequestDescription             string     `json:"request_description"`
	FRTQuantity                    float64    `json:"frt_quantity"`
	OperationItemPrice             float64    `json:"operation_item_price"`
	DiscountItemAmount             float64    `json:"discount_item_amount"`
	DiscountItemPercent            float64    `json:"discount_item_percent"`
	DiscountRequestPercent         float64    `json:"discount_request_percent"`
	DiscountRequestAmount          float64    `json:"discount_request_amount"`
	DiscountApprovalBy             string     `json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `json:"discount_approval_date"`
}

type BookEstimationPayloadsDiscount struct{
	PackageDiscount int `json:"disccount_request_percent"`
	Operation int `json:"disccount_request_percent"`
	Sparepart int `json:"disccount_request_percent"`
	Oil int `json:"disccount_request_percent"`
	Material int `json:"disccount_request_percent"`
	Fee int `json:"disccount_request_percent"`
	Accessories int `json:"disccount_request_percent"`
	Souvenir int `json:"disccount_request_percent"`
}

type BookingEstimationPostPayloads struct{
	BrandId int `json:"brand_id"`
	ModelId int `json:"model_id"`
	VehicleId int `json:"vehicle_id"`
	DealerRepCodeId int `json:"dealer_rep_code_id"`
	ContactPersonName string `json:"contact_person_name"`
	ContactPersonPhone string `json:"contact_person_phne"`
	ContactPersonMobile string `json:"contact_person_mobile"`
	ContactViaId int `json:"contact_person_via"`
	InsPolicyNo string `json:"insurance_policy_number"`  
	InsExpireDate time.Time `json:"insurance_expire_date"`
	InsClaimNo string `json:"insurance_claim_no"`
	InsPIC string `json:"insurace_pic"`
	CampagnCodeId int `json:"campaign_code_id"`
}