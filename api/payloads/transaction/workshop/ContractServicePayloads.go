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
	VehicleId                     int       `gorm:"column:vehicle_id" json:"vehicle_id" parent_entity:"trx_contract_service`
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
		VehicleId             int    `json:"vehicle_id"`
		IsActive              bool   `json:"is_active"`
		VehicleChassisNumber  string `json:"vehicle_chassis_number"`
		VehicleEngineNumber   string `json:"vehicle_engine_number"`
		VehicleBrandId        int    `json:"vehicle_brand_id"`
		VehicleModelId        int    `json:"vehicle_model_id"`
		VehicleVariantId      int    `json:"vehicle_variant_id"`
		VehicleColourId       int    `json:"vehicle_colour_id"`
		OptionId              int    `json:"option_id"`
		VehicleTransmissionId int    `json:"vehicle_transmission_id"`
		VehicleWheelDriveId   int    `json:"vehicle_wheel_drive_id"`
		VehicleCylinderId     int    `json:"vehicle_cylinder_id"`
		VehicleFuelId         int    `json:"vehicle_fuel_id"`
		//VehicleProductionYear       string      `json:"vehicle_production_year"`
		//VehicleAssemblyYear         int         `json:"vehicle_assembly_year"`
		VehicleVin2                 string      `json:"vehicle_vin2"`
		VehicleKeyCode              string      `json:"vehicle_key_code"`
		VehicleRadiotapeCode        string      `json:"vehicle_radiotape_code"`
		VehicleServiceBookingNumber interface{} `json:"vehicle_service_booking_number"`
		VehicleHandoverDocumentDate string      `json:"vehicle_handover_document_date"`
		VehicleDealerStock          string      `json:"vehicle_dealer_stock"`
		//VehicleLastKm               int         `json:"vehicle_last_km"`
		VehicleStatusId             int    `json:"vehicle_status_id"`
		VehicleIsGreyMarket         bool   `json:"vehicle_is_grey_market"`
		AccessoriesOptionCode       string `json:"accessories_option_code"`
		AccessoriesOptionName       string `json:"accessories_option_name"`
		VehicleVariantDescription   string `json:"vehicle_variant_description"`
		VehicleVariantCode          string `json:"vehicle_variant_code"`
		VehicleLastServiceDate      string `json:"vehicle_last_service_date"`
		VehicleColourCode           string `json:"vehicle_colour_code"`
		VehicleColourCommercialName string `json:"vehicle_colour_commercial_name"`
	} `json:"master"`
	Stnk struct {
		VehicleRegistrationCertificateTnkb      string `json:"vehicle_registration_certificate_tnkb"`
		VehicleRegistrationCertificateNumber    string `json:"vehicle_registration_certificate_number"`
		VehicleRegistrationCertificateValidDate string `json:"vehicle_registration_certificate_valid_date"`
		VehicleRegistrationCertificateCode      string `json:"vehicle_registration_certificate_code"`
		VehicleRegistrationCertificateOwnerName string `json:"vehicle_registration_certificate_owner_name"`
		DetailAddress                           struct {
			StnkAddressId      int    `json:"stnk_address_id"`
			StnkAddressStreet1 string `json:"stnk_address_street_1"`
			StnkAddressStreet2 string `json:"stnk_address_street_2"`
			StnkAddressStreet3 string `json:"stnk_address_street_3"`
			StnkVillageCode    string `json:"stnk_village_code"`
			StnkVillageDesc    string `json:"stnk_village_desc"`
			StnkSubdistrict    string `json:"stnk_subdistrict"`
			StnkProvince       string `json:"stnk_province"`
			StnkCity           string `json:"stnk_city"`
			StnkZipCode        string `json:"stnk_zip_code"`
		} `json:"detail_address"`
		StnkVillageId                            int         `json:"stnk_village_id"`
		VehicleRegistrationCertificateTnkbColour string      `json:"vehicle_registration_certificate_tnkb_colour"`
		VehicleBpkbNumber                        string      `json:"vehicle_bpkb_number"`
		VehicleBpkbDate                          string      `json:"vehicle_bpkb_date"`
		NikCertificationNumber                   interface{} `json:"nik_certification_number"`
		PoliceInvoiceNumber                      string      `json:"police_invoice_number"`
	} `json:"stnk"`
	//CurrentUser struct {
	//	UserCustomerId    int    `json:"user_customer_id"`
	//	UserCustomerCode  string `json:"user_customer_code"`
	//	UserCustomerName  string `json:"user_customer_name"`
	//	DetailUserAddress struct {
	//		UserAddressId      int    `json:"user_address_id"`
	//		UserAddressStreet1 string `json:"user_address_street_1"`
	//		UserAddressStreet2 string `json:"user_address_street_2"`
	//		UserAddressStreet3 string `json:"user_address_street_3"`
	//		UserVillageCode    string `json:"user_village_code"`
	//		UserVillageDesc    string `json:"user_village_desc"`
	//		UserSubdistrict    string `json:"user_subdistrict"`
	//		UserProvince       string `json:"user_province"`
	//		UserCity           string `json:"user_city"`
	//		UserZipCode        string `json:"user_zip_code"`
	//		UserNpwpNo         string `json:"user_npwp_no"`
	//	} `json:"detail_user_address"`
	//} `json:"current_user"`
	//Billing struct {
	//	BillingCustomerId    int    `json:"billing_customer_id"`
	//	BillingCustomerCode  string `json:"billing_customer_code"`
	//	BillingCustomerName  string `json:"billing_customer_name"`
	//	DetailBillingAddress struct {
	//		BillingAddressId      int    `json:"billing_address_id"`
	//		BillingAddressStreet1 string `json:"billing_address_street_1"`
	//		BillingAddressStreet2 string `json:"billing_address_street_2"`
	//		BillingAddressStreet3 string `json:"billing_address_street_3"`
	//		BillingVillageCode    string `json:"billing_village_code"`
	//		BillingVillageDesc    string `json:"billing_village_desc"`
	//		BillingSubdistrict    string `json:"billing_subdistrict"`
	//		BillingProvince       string `json:"billing_province"`
	//		BillingCity           string `json:"billing_city"`
	//		BillingZipCode        string `json:"billing_zip_code"`
	//	} `json:"detail_billing_address"`
	//} `json:"billing"`
	//ContractService struct {
	//	//ContractId            int    `json:"contract_id"`
	//	ContractServiceNumber string `json:"contract_service_number"`
	//	ContractServiceDate   string `json:"contract_service_date"`
	//	ContractDealer        int    `json:"contract_dealer"`
	//} `json:"contract_service"`
	//Insurance struct {
	//	InsuranceCompany      string `json:"insurance_company"`
	//	InsurancePolicyNumber string `json:"insurance_policy_number"`
	//	InsuranceEndDate      string `json:"insurance_end_date"`
	//	InsuranceEraNumber    string `json:"insurance_era_number"`
	//} `json:"insurance"`
	//Receiving struct {
	//	VehicleWrsNumber           string      `json:"vehicle_wrs_number"`
	//	VehicleWrsDate             string      `json:"vehicle_wrs_date"`
	//	VehicleBpuNumber           interface{} `json:"vehicle_bpu_number"`
	//	VehicleBpuDate             interface{} `json:"vehicle_bpu_date"`
	//	VehicleDeliveryOrderNumber interface{} `json:"vehicle_delivery_order_number"`
	//	VehicleDeliveryOrderDate   interface{} `json:"vehicle_delivery_order_date"`
	//	VehicleSjNumber            string      `json:"vehicle_sj_number"`
	//	VehicleSjDate              string      `json:"vehicle_sj_date"`
	//	VehicleExpedition          interface{} `json:"vehicle_expedition"`
	//} `json:"receiving"`
	//LastStatus struct {
	//	BookedBy                        string      `json:"booked_by"`
	//	BookingNumber                   string      `json:"booking_number"`
	//	BookingDate                     string      `json:"booking_date"`
	//	VehicleSalesOrderDocumentNumber interface{} `json:"vehicle_sales_order_document_number"`
	//	VehicleSalesOrderDate           string      `json:"vehicle_sales_order_date"`
	//	VehicleLastFakturPolisiNumber   string      `json:"vehicle_last_faktur_polisi_number"`
	//	VehicleLastStnkNumber           string      `json:"vehicle_last_stnk_number"`
	//	VehicleLastStnkDate             string      `json:"vehicle_last_stnk_date"`
	//	VehicleLastBpkbNumber           string      `json:"vehicle_last_bpkb_number"`
	//	VehicleLastBpkbDate             string      `json:"vehicle_last_bpkb_date"`
	//	VehicleLastNikNumber            string      `json:"vehicle_last_nik_number"`
	//	BpkDocumentNumber               interface{} `json:"bpk_document_number"`
	//	BpkDate                         string      `json:"bpk_date"`
	//	VehicleUnitDelivery             string      `json:"vehicle_unit_delivery"`
	//	VehicleUnitStockStatus          string      `json:"vehicle_unit_stock_status"`
	//	WarehouseId                     int         `json:"warehouse_id"`
	//	WarehouseName                   string      `json:"warehouse_name"`
	//	VehicleUnitLotNumber            string      `json:"vehicle_unit_lot_number"`
	//} `json:"last_status"`
	//Purchase struct {
	//	VehiclePurchaseDate    string `json:"vehicle_purchase_date"`
	//	VehiclePurchaseOrderId int    `json:"vehicle_purchase_order_id"`
	//	VehiclePurchaseDealer  int    `json:"vehicle_purchase_dealer"`
	//	VehicleCustomerCode    string `json:"vehicle_customer_code"`
	//	VehiclePurchasePrice   int    `json:"vehicle_purchase_price"`
	//} `json:"purchase"`
	////VehicleId           int    `json:"vehicle_id"`
	////VehicleCode         string `json:"vehicle_chassis_number"`
	////VehicleEngineNumber string `json:"vehicle_engine_number"`
	////VehicleTnkb         string `json:"vehicle_registration_certificate_tnkb"`
	////VehicleBrandId      int    `json:"vehicle_brand_id"`
	////VehicleModelId      int    `json:"vehicle_model_id"`
	////VehicleModelVariant string `json:"model_variant_colour_description"`
	////VehicleVariantId    int    `json:"vehicle_variant_id"`
	////VehicleColourId     int    `json:"vehicle_colour_id"`
	//VehicleOwner        string `json:"vehicle_registration_certificate_owner_name"`
}

type SubmitContractServiceResponse struct {
	ContractServiceSystemNumber   int    `json:"contract_service_system_number"`
	ContractSevriceDocumentNumber string `json:"contract_service_document_number"`
}

// type ContractServiceColour struct {
// 	VariantColourId   int    `json:"colour_id"`
// 	VariantColourCode string `json:"colour_commercial_name"`
// 	VariantColourName string `json:"colour_police_name"`
// }
