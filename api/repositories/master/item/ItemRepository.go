package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]masteritempayloads.ItemLookup, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLookup(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemCode(*gorm.DB, string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveItem(*gorm.DB, masteritempayloads.ItemResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
