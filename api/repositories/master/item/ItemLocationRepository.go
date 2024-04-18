package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLocationRepository interface {
	AddItemLocation(*gorm.DB, masteritempayloads.ItemLocationDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveItemLocation(*gorm.DB, masteritempayloads.ItemLocationRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetItemLocationById(tx *gorm.DB, Id int) (masteritempayloads.ItemLocationRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
