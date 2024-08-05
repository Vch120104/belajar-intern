package transactionunitpayloads

type PdiRequest struct {
	BrandID              int     `json:"brand_id"`
	PdiDocumentNumber    string  `json:"pdi_document_number"`
	ModelID              int     `json:"model_id"`
	VariantID            int     `json:"variant_id"`
	VehicleID            int     `json:"vehicle_id"`
	CompanyID            int     `json:"company_id"`
	OperationNumber      string  `json:"operation_number"`
	Frt                  float64 `json:"frt"`
	ContractSystemNumber int     `json:"contract_system_number"`
}
