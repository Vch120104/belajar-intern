package transactionworkshoppayloads

import "time"

type ContractServiceResponse struct {
	CompanyId                     int       `json:"company_id"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number"`
	ContractServiceFrom           time.Time `json:"contract_service_from"`
	ContractServiceTo             time.Time `json:"contract_service_to"`
	BrandId                       int       `json:"brand_id"`
	BrandCode                     string    `json:"brand_code"`
	BrandName                     string    `json:"brand_name"`
	ModelId                       int       `json:"model_id"`
	ModelCode                     string    `json:"model_code"`
	ModelName                     string    `json:"model_description"`
	ModelCodeDescription          string    `json:"model_code_description"`
	VehicleId                     int       `json:"vehicle_id"`
	VehicleCode                   string    `json:"vehicle_chassis_number"`
	VehicleEngineNumber           string    `json:"vehicle_engine_number"`
	VehicleTnkb                   string    `json:"vehicle_registration_certificate_tnkb"`
	ContractServiceStatusId       int       `json:"contract_service_status_id"`
}

type ContractServiceResponseId struct {
	CompanyId                     int       `json:"company_id"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number"`
	ContractServiceFrom           time.Time `json:"contract_service_from"`
	ContractServiceTo             time.Time `json:"contract_service_to"`
	BrandId                       int       `json:"brand_id"`
	BrandCode                     string    `json:"brand_code"`
	BrandName                     string    `json:"brand_name"`
	ModelId                       int       `json:"model_id"`
	ModelCode                     string    `json:"model_code"`
	ModelName                     string    `json:"model_description"`
	ModelCodeDescription          string    `json:"model_code_description"`
	VehicleId                     int       `json:"vehicle_id"`
	VehicleCode                   string    `json:"vehicle_chassis_number"`
	VehicleEngineNumber           string    `json:"vehicle_engine_number"`
	VehicleTnkb                   string    `json:"vehicle_registration_certificate_tnkb"`
	ContractServiceStatusId       int       `json:"contract_service_status_id"`
}

type ContractServiceRequest struct {
	CompanyId                     int       `json:"company_id" parent_entity:"trx_contract_service" main_table:"trx_contract_service"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number" parent_entity:"trx_contract_service"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number" parent_entity:"trx_contract_service"`
	ContractServiceFrom           time.Time `json:"contract_service_from" parent_entity:"trx_contract_service"`
	ContractServiceTo             time.Time `json:"contract_service_to" parent_entity:"trx_contract_service"`
	BrandId                       int       `json:"brand_id" parent_entity:"trx_contract_service"`
	ModelId                       int       `json:"model_id" parent_entity:"trx_contract_service"`
	// ModelCodeDescription          string    `json:"model_code_description" parent_entity:"trx_contract_service"`
	VehicleId               int `json:"vehicle_id" parent_entity:"trx_contract_service"`
	ContractServiceStatusId int `json:"contract_service_status_id" parent_entity:"trx_contract_service"`
}

type ContractServiceBrand struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type ContractServiceModel struct {
	ModelId              int    `json:"model_id"`
	ModelCode            string `json:"model_code"`
	ModelName            string `json:"model_description"`
	ModelCodeDescription string `json:"model_code_description"`
}

type ContractServiceVariant struct {
	VariantId   int    `json:"variant_id"`
	VariantCode string `json:"variant_code"`
	VariantName string `json:"variant_description"`
}

type ContractServiceVehicleResponse struct {
	VehicleId           int    `json:"vehicle_id"`
	VehicleCode         string `json:"vehicle_chassis_number"`
	VehicleEngineNumber string `json:"vehicle_engine_number"`
	VehicleTnkb         string `json:"vehicle_registration_certificate_tnkb"`
	// VehicleProduction   json.RawMessage `json:"vehicle_production_year"`
	// VehicleLastKm       json.RawMessage `json:"vehicle_last_km"`
	VehicleBrandId      int    `json:"vehicle_brand_id"`
	VehicleModelId      int    `json:"vehicle_model_id"`
	VehicleModelVariant string `json:"model_variant_colour_description"`
	VehicleVariantId    int    `json:"vehicle_variant_id"`
	VehicleColourId     int    `json:"vehicle_colour_id"`
	VehicleOwner        string `json:"vehicle_registration_certificate_owner_name"`
}

// type ContractServiceColour struct {
// 	VariantColourId   int    `json:"colour_id"`
// 	VariantColourCode string `json:"colour_commercial_name"`
// 	VariantColourName string `json:"colour_police_name"`
// }
