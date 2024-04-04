package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarrantyFreeServiceRepository interface {
	GetAllWarrantyFreeService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusWarrantyFreeService(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
