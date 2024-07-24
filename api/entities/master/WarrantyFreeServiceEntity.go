package masterentities

import "time"

var CreateWarrantyFreeServiceTable = "mtr_warranty_free_service"

type WarrantyFreeService struct {
	IsActive                      bool      `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	WarrantyFreeServicesId        int       `gorm:"column:warranty_free_services_id;size:30;not null;primaryKey"        json:"warranty_free_services_id"`
	BrandId                       int       `gorm:"column:brand_id;size:30;not null"        json:"brand_id"`                                           //Fk with mtr_brand on sales service
	ModelId                       int       `gorm:"column:model_id;size:30;not null"        json:"model_id"`                                           //Fk with mtr_unit_model on sales service
	WarrantyFreeServiceTypeId     int       `gorm:"column:warranty_free_service_type_id;size:30;not null"        json:"warranty_free_service_type_id"` //Fk with mtr_warranty_free_service_type on general service
	EffectiveDate                 time.Time `gorm:"column:effective_date;not null"        json:"effective_date"`
	ExpireMileage                 float64   `gorm:"column:expire_mileage;not null"        json:"expire_mileage"`
	ExpireMonth                   float64   `gorm:"column:expire_month;not null"        json:"expire_month"`
	VariantId                     int       `gorm:"column:variant_id;size:30;null"        json:"variant_id"` //Fk with mtr_unit_variant on sales service
	ExpireMileageExtendedWarranty float64   `gorm:"column:expire_mileage_extended_warranty;null"        json:"expire_mileage_extended_warranty"`
	ExpireMonthExtendedWarranty   float64   `gorm:"column:expire_month_extended_warranty;null"        json:"expire_month_extended_warranty"`
	Remark                        string    `gorm:"column:remark;size:256;null"        json:"remark"`
	ExtendedWarranty              *bool      `gorm:"column:extended_warranty;null"        json:"extended_warranty"`
}

func (*WarrantyFreeService) TableName() string {
	return CreateWarrantyFreeServiceTable
}
