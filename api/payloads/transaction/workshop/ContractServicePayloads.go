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
	VehicleOwner                  string    `json:"vehicle_registration_certificate_owner_name"`
}

type ContractServiceRequest struct {
	CompanyId                     int       `json:"company_id" parent_entity:"trx_contract_service" main_table:"trx_contract_service"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number" parent_entity:"trx_contract_service"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number" parent_entity:"trx_contract_service"`
	ContractServiceFrom           time.Time `json:"contract_service_from" parent_entity:"trx_contract_service"`
	ContractServiceTo             time.Time `json:"contract_service_to" parent_entity:"trx_contract_service"`
	ContractServiceDate           time.Time `json:"contract_service_date" parent_entity:"trx_contract_service"`
	BrandId                       int       `json:"brand_id" parent_entity:"trx_contract_service"`
	ModelId                       int       `json:"model_id" parent_entity:"trx_contract_service"`
	VehicleId                     int       `json:"vehicle_id" parent_entity:"trx_contract_service"`
	ContractServiceStatusId       int       `json:"contract_service_status_id" parent_entity:"trx_contract_service"`
}

type ContractServiceInsert struct {
	CompanyId                     int       `json:"company_id"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number"`
	ContractServiceDate           time.Time `json:"contract_service_date"`
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
	VehicleOwner                  string    `json:"vehicle_registration_certificate_owner_name"`
	RegisteredMileage             int       `json:"registered_mileage"`
	Remark                        string    `json:"remark"`
	Total                         float64   `json:"total" default:"0"`
	Vat                           float64   `json:"vat" default:"0"`
	GrandTotal                    float64   `json:"grand_total" default:"0"`
}

type ContractServiceInsertResponse struct {
	CompanyId                     int       `json:"company_id"`
	ContractServiceSystemNumber   int       `json:"contract_service_system_number"`
	ContractServiceDocumentNumber string    `json:"contract_service_document_number"`
	ContractServiceDate           time.Time `json:"contract_service_date"`
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
	VehicleOwner                  string    `json:"vehicle_registration_certificate_owner_name"`
	RegisteredMileage             int       `json:"registered_mileage"`
	Remark                        string    `json:"remark"`
	Total                         float64   `json:"total" default:"0"`
	Vat                           float64   `json:"vat" default:"0"`
	GrandTotal                    float64   `json:"grand_total" default:"0"`
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
	Master struct {
		VehicleId                   int     `json:"vehicle_id"`
		VehicleCode                 string  `json:"vehicle_chassis_number"`
		VehicleEngineNumber         string  `json:"vehicle_engine_number"`
		VehicleBrandId              int     `json:"vehicle_brand_id"`
		VehicleModelId              int     `json:"vehicle_model_id"`
		VehicleVariantId            int     `json:"vehicle_variant_id"`
		VehicleColourId             int     `json:"vehicle_colour_id"`
		VehicleTransmissionId       int     `json:"vehicle_transmission_id"`
		VehicleWheelDriveId         int     `json:"vehicle_wheel_drive_id"`
		VehicleCylinderId           int     `json:"vehicle_cylinder_id"`
		VehicleFuelId               int     `json:"vehicle_fuel_id"`
		VehicleProductionYear       float64 `json:"vehicle_production_year"`
		VehicleAssemblyYear         float64 `json:"vehicle_assembly_year"`
		VehicleStatusId             int     `json:"vehicle_status_id"`
		VehicleIsGreyMarket         bool    `json:"vehicle_is_grey_market"`
		VehicleVariantDescription   string  `json:"vehicle_variant_description"`
		VehicleVariantCode          string  `json:"vehicle_variant_code"`
		VehicleColourCode           string  `json:"vehicle_colour_code"`
		VehicleColourCommercialName string  `json:"vehicle_colour_commercial_name"`
	} `json:"master"`

	STNK ContractServiceVehicleSTNKResponse `json:"stnk"`
}

type ContractServiceVehicleSTNKResponse struct {
	VehicleRegistrationCertificateTnkb      string `json:"vehicle_registration_certificate_tnkb"`
	VehicleRegistrationCertificateNumber    string `json:"vehicle_registration_certificate_number"`
	VehicleRegistrationCertificateValidDate string `json:"vehicle_registration_certificate_valid_date"`
	VehicleRegistrationCertificateCode      string `json:"vehicle_registration_certificate_code"`
	VehicleRegistrationCertificateOwnerName string `json:"vehicle_registration_certificate_owner_name"`
}

type SubmitContractServiceResponse struct {
	ContractServiceSystemNumber   int    `json:"contract_service_system_number"`
	ContractSevriceDocumentNumber string `json:"contract_service_document_number"`
}
