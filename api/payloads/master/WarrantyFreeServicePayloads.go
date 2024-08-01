package masterpayloads

import "time"

type WarrantyFreeServiceResponse struct {
	IsActive                      bool    `json:"is_active"`
	WarrantyFreeServicesId        int     `json:"warranty_free_services_id"`
	BrandId                       int     `json:"brand_id"`
	ModelId                       int     `json:"model_id"`
	WarrantyFreeServiceTypeId     int     `json:"warranty_free_service_type_id"`
	EffectiveDate                 string  `json:"effective_date"`
	ExpireMileage                 float64 `json:"expire_mileage"`
	ExpireMonth                   float64 `json:"expire_month"`
	VariantId                     int     `json:"variant_id"`
	ExpireMileageExtendedWarranty float64 `json:"expire_mileage_extended_warranty"`
	ExpireMonthExtendedWarranty   float64 `json:"expire_month_extended_warranty"`
	Remark                        string  `json:"remark"`
	ExtendedWarranty              bool    `json:"extended_warranty"`
}

type WarrantyFreeServiceListResponse struct {
	IsActive                  bool   `json:"is_active" parent_entity:"mtr_warranty_free_service"`
	WarrantyFreeServicesId    int    `json:"warranty_free_services_id" parent_entity:"mtr_warranty_free_service" main_table:"mtr_warranty_free_service"`
	BrandId                   int    `json:"brand_id" parent_entity:"mtr_warranty_free_service"`
	ModelId                   int    `json:"model_id" parent_entity:"mtr_warranty_free_service"`
	WarrantyFreeServiceTypeId int    `json:"warranty_free_service_type_id" parent_entity:"mtr_warranty_free_service"`
	EffectiveDate             string `json:"effective_date" parent_entity:"mtr_warranty_free_service"`
}

type WarrantyFreeServiceRequest struct {
	WarrantyFreeServicesId        int       `json:"warranty_free_services_id"`
	BrandId                       int       `json:"brand_id"`
	ModelId                       int       `json:"model_id"`
	WarrantyFreeServiceTypeId     int       `json:"warranty_free_service_type_id"`
	EffectiveDate                 time.Time `json:"effective_date"`
	ExpireMileage                 float64   `json:"expire_mileage"`
	ExpireMonth                   float64   `json:"expire_month"`
	VariantId                     int       `json:"variant_id"`
	ExpireMileageExtendedWarranty float64   `json:"expire_mileage_extended_warranty"` //default value 0
	ExpireMonthExtendedWarranty   float64   `json:"expire_month_extended_warranty"`   //default value 0
	Remark                        string    `json:"remark"`                           // default value ""
	ExtendedWarranty              bool      `json:"extended_warranty"`
}

type WarrantyFreeServicePatchResponse struct {
	IsActive                      bool    `json:"is_active"`
	WarrantyFreeServicesId        int     `json:"warranty_free_services_id"`
}

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandName string `json:"brand_name"`
	BrandCode string `json:"brand_code"`
}

type UnitModelResponse struct {
	ModelId          int    `json:"model_id"`
	ModelDescription string `json:"model_description"`
	ModelCode        string `json:"model_code"`
}

type UnitVariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantDescription string `json:"variant_description"`
}

type WarrantyFreeServiceTypeResponse struct {
	WarrantyFreeServiceTypeId          int    `json:"warranty_free_service_type_id"`
	WarrantyFreeServiceTypeCode        string `json:"warranty_free_service_type_code"`
	WarrantyFreeServiceTypeDescription string `json:"warranty_free_service_type_description"`
}
