package transactionworkshoppayloads

import "time"

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
