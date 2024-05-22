package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarrantyFreeServiceRepository interface {
	GetAllWarrantyFreeService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusWarrantyFreeService(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
