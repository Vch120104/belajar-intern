package transactionworkshoppayloads

import "time"

type AtpmClaimRegistrationRequest struct {
	ClaimSystemNumber     int       `json:"claim_system_number"`
	CompanyId             int       `json:"company_id"`
	BrandId               int       `json:"brand_id"`
	ModelId               int       `json:"model_id"`
	VariantId             int       `json:"variant_id"`
	FspCategoryId         int       `json:"fsp_category_id"`
	ClaimTypeId           int       `json:"claim_type_id"`
	WarrantyTypeId        int       `json:"warranty_type_id"`
	CustomerComplaint     string    `json:"customer_complaint"`
	TechnicianDiagnostic  string    `json:"technician_diagnostic"`
	Countermeasure        string    `json:"countermeasure"`
	ClaimDate             time.Time `json:"claim_date"`
	RepairEndDate         time.Time `json:"repair_end_date"`
	Fuel                  float64   `json:"fuel"`
	CustomerId            int       `json:"customer_id"`
	VDN                   string    `json:"vdn"`
	ClaimHeader           string    `json:"claim_header"`
	WorkOrderSystemNumber int       `json:"work_order_system_number"`
}

type AtpmClaimRegistrationRequestSave struct {
	CompanyId            int       `json:"company_id"`
	CustomerComplaint    string    `json:"customer_complaint"`
	TechnicianDiagnostic string    `json:"technician_diagnostic"`
	Countermeasure       string    `json:"countermeasure"`
	RepairEndDate        time.Time `json:"repair_end_date"`
	Fuel                 float64   `json:"fuel"`
	CustomerId           int       `json:"customer_id"`
	VDN                  string    `json:"vdn"`
	ClaimHeader          string    `json:"claim_header"`
	SymptomCode          string    `json:"symptom_code"`
	TroubleCode          string    `json:"trouble_code"`
}

type AtpmClaimRegistrationResponse struct {
	ClaimSystemNumber       int       `json:"claim_system_number"`
	CompanyId               int       `json:"company_id"`
	CompanyName             string    `json:"company_name"`
	BrandId                 int       `json:"brand_id"`
	BrandName               string    `json:"brand_name"`
	ClaimTo                 string    `json:"claim_to"`
	ClaimTypeId             int       `json:"claim_type_id"`
	ClaimTypeName           string    `json:"claim_type_name"`
	ClaimNumber             string    `json:"claim_number"`
	ClaimDate               time.Time `json:"claim_date"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number"`
	WorkOrderSystemNumber   string    `json:"work_order_system_number"`
	WorkOrderDate           time.Time `json:"work_order_date"`
	QCPassDate              time.Time `json:"qc_pass_date"`
	RetailSalesDate         time.Time `json:"retail_sales_date"`
	RepairEndDate           time.Time `json:"repair_end_date"`
	VehicleId               int       `json:"vehicle_id"`
	VehicleChassisNumber    string    `json:"vehicle_chassis_number"`
	VehicleEngineNumber     string    `json:"vehicle_engine_number"`
	VehicleTnkb             string    `json:"vehicle_tnkb"`
	ModelId                 int       `json:"model_id"`
	ModelDescription        string    `json:"model_description"`
	VariantId               int       `json:"variant_id"`
	VariantDescription      string    `json:"variant_description"`
	Km                      string    `json:"km"`
	ServiceBook             string    `json:"service_book"`
	DeliveryDate            time.Time `json:"delivery_date"`
}

type AtpmClaimRegistrationDetailResponse struct {
	ClaimDetailSystemNumber  int     `json:"claim_detail_system_number"`
	ClaimSystemNumber        int     `json:"claim_system_number"`
	CompanyId                int     `json:"company_id"`
	ClaimLineNumber          int     `json:"claim_line_number"`
	LineTypeId               int     `json:"line_type_id"`
	LineTypeName             string  `json:"line_type_name"`
	ItemOperationId          int     `json:"item_operation_id"`
	ItemOperationCode        string  `json:"item_operation_code"`
	ItemOperationDescription string  `json:"item_operation_description"`
	FrtQuantity              float64 `json:"frt_quantity"`
	Uom                      string  `json:"uom"`
	ItemPrice                float64 `json:"item_price"`
	DiscountPercent          float64 `json:"discount_percent"`
	TotalAmount              float64 `json:"total_amount"`
	RecallNumber             int     `json:"recall_number"`
	ReturnToAtpm             bool    `json:"return_to_atpm"`
}

type AtpmClaimDetailRequest struct {
	ClaimSystemNumber     int     `json:"claim_system_number"`
	CompanyId             int     `json:"company_id"`
	ClaimLineNumber       int     `json:"claim_line_number"`
	WorkOrderSystemNumber int     `json:"work_order_system_number"`
	WorkOrderLineNumber   int     `json:"work_order_line_number"`
	LineTypeId            int     `json:"line_type_id"`
	ItemOperationId       int     `json:"item_operation_id"`
	FrtQuantity           float64 `json:"frt_quantity"`
	DiscountPercent       float64 `json:"discount_percent"`
	DiscountAmount        float64 `json:"discount_amount"`
}
