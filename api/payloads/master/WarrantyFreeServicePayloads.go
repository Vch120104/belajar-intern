package masterpayloads

import "time"

type WarrantyFreeServiceResponse struct {
	IsActive                  bool `json:"is_active"`
	WarrantyFreeServicesId    int  `json:"warranty_free_services_id"`
	BrandId                   int  `json:"brand_id"`
	ModelId                   int  `json:"model_id"`
	WarrantyFreeServiceTypeId int  `json:"warranty_free_service_type_id"`
	// EffectiveDate                 time.Time `json:"effective_date"`
	ExpireMileage                 float64 `json:"expire_mileage"`
	ExpireMonth                   float64 `json:"expire_month"`
	VariantId                     int     `json:"variant_id"`
	ExpireMileageExtendedWarranty float64 `json:"expire_mileage_extended_warranty"`
	ExpireMonthExtendedWarranty   float64 `json:"expire_month_extended_warranty"`
	Remark                        string  `json:"remark"`
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
}

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandName string `json:"brand_name"`
}

type UnitModelResponse struct {
	ModelId          int    `json:"model_id"`
	ModelDescription string `json:"model_description"`
}

type UnitVariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantDescription string `json:"variant_description"`
}