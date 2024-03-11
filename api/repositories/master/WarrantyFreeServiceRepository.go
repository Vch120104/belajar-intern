package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type WarrantyFreeServiceRepository interface {
	GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, error)
	SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (bool, error)
}
