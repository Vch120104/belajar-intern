package masteroperationpayloads

type OperationModelMappingResponse struct {
	IsActive                bool `json:"is_active"`
	OperationModelMappingId int  `json:"operation_model_mapping_id"`
	BrandId                 int  `json:"brand_id"`
	ModelId                 int  `json:"model_id"`
	OperationId             int  `json:"operation_id"`
	OperationUsingIncentive bool `json:"operation_using_incentive"`
	OperationUsingActual    bool `json:"operation_using_actual"`
	OperationPdi            bool `json:"operation_pdi"`
}

type OperationModelMappingFrtRequest struct {
	IsActive                bool    `json:"is_active"`
	OperationModelMappingId int     `json:"operation_model_mapping_id"`
	OperationFrtId          int     `json:"operation_frt_id"`
	VariantId               int     `json:"variant_id"`
	FrtHour                 float64 `json:"frt_hour"`
	FrtHourExpress          float64 `json:"frt_hour_express"`
}

type OperationModelMappingRequest struct {
	BrandId                 int  `json:"brand_id"`
	ModelId                 int  `json:"model_id"`
	OperationId             int  `json:"operation_id"`
	OperationUsingIncentive bool `json:"operation_using_incentive"`
	OperationUsingActual    bool `json:"operation_using_actual"`
	OperationPdi            bool `json:"operation_pdi"`
}

type OperationModelMappingUpdate struct {
	OperationModelMappingId int  `json:"operation_model_mapping_id"`
	OperationUsingIncentive bool `json:"operation_using_incentive"`
	OperationUsingActual    bool `json:"operation_using_actual"`
	OperationPdi            bool `json:"operation_pdi"`
}

type HeaderRequest struct {
	BrandId                 int  `json:"brand_id"`
	ModelId                 int  `json:"model_id"`
	OperationId             int  `json:"operation_id"`
	OperationUsingIncentive bool `json:"operation_using_incentive"`
	OperationUsingActual    bool `json:"operation_using_actual"`
	OperationPdi            bool `json:"operation_pdi"`
}

type DetailRequest struct {
	IsActive                bool    `json:"is_active"`
	OperationModelMappingId int     `json:"operation_model_mapping_id"`
	OperationFrtId          int     `json:"operation_frt_id"`
	VariantId               int     `json:"variant_id"`
	FrtHour                 float64 `json:"frt_hour"`
	FrtHourExpress          float64 `json:"frt_hour_express"`
}

type OperationModelMappingAndFRTRequest struct {
	HeaderRequest OperationModelMappingResponse   `json:"headerRequest"`
	DetailRequest OperationModelMappingFrtRequest `json:"detailRequest"`
}

type OperationModelMappingDocumentRequirementRequest struct {
	IsActive                                bool   `json:"is_active"`
	OperationModelMappingId                 int    `json:"operation_model_mapping_id"`
	OperationDocumentRequirementId          int    `json:"operation_document_requirement_id"`
	Line                                    int    `json:"line"`
	OperationDocumentRequirementDescription string `json:"operation_document_requirement_description"`
}

type OperationModelMappingLookup struct {
	IsActive                bool   `json:"is_active" parent_entity:"mtr_operation_model_mapping"`
	OperationModelMappingId int    `json:"operation_model_mapping_id" parent_entity:"mtr_operation_model_mapping" main_table:"mtr_operation_model_mapping"`
	OperationId             int    `json:"operation_id" parent_entity:"mtr_operation_code" references:"mtr_operation_code"`
	OperationName           string `json:"operation_name" parent_entity:"mtr_operation_code"`
	OperationCode           string `json:"operation_code" parent_entity:"mtr_operation_code"`
	BrandId                 int    `json:"brand_id" parent_entity:"mtr_operation_model_mapping"`
	ModelId                 int    `json:"model_id" parent_entity:"mtr_operation_model_mapping"`
}

type OperationModelModelBrandOperationCodeRequest struct {
	BrandId     int `json:"brand_id"`
	ModelId     int `json:"model_id"`
	OperationId int `json:"operation_id"`
}

// SELECT
// [Operation Code] = A.OPERATION_CODE,
// [Operation Description] = B.OPERATION_NAME,
// [Vehicle Brand] = A.VEHICLE_BRAND,
// [Model Code] = A.MODEL_CODE
// FROM amOperation0 A
// LEFT JOIN amOperationCode B ON A.OPERATION_CODE = B.OPERATION_CODE

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type VariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
}

type CurrencyResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
}

type CompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
}

type ModelResponse struct {
	ModelId          int    `json:"model_id"`
	ModelDescription string `json:"model_description"`
}

type OperationLevelByIdResponse struct {
	IsActive                    bool   `json:"is_active"`
	OperationLevelId            int    `json:"operation_level_id"`
	OperationModelMappingId     int    `json:"operation_model_mapping_id"`
	OperationEntriesId          int    `json:"operation_entries_id"`
	OperationEntriesCode        string `json:"operation_entries_code"`
	OperationEntriesDescription string `json:"operation_entries_description"`
	OperationGroupCode          int    `json:"operation_group_code"`
	OperationGroupDescription   string `json:"operation_group_description"`
	OperationKeyCode            string `json:"operation_key_code"`
	OperationKeyDescription     string `json:"operation_key_description"`
	OperationSectionCode        string `json:"operation_section_code"`
	OperationSectionDescription string `json:"operation_section_description"`
}
