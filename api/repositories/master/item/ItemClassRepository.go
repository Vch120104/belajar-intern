package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassRepository interface {
	GetAllItemClass(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetItemClassById(tx *gorm.DB, Id int) (masteritempayloads.ItemClassResponse, *exceptionsss_test.BaseErrorResponse)
	SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemClass(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
