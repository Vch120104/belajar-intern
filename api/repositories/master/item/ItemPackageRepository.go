package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepository interface {
	GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetItemPackageById(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
