package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, paginate pagination.Pagination) ([]map[string]interface{}, int, int, error)
	GetAllItemLookup(tx *gorm.DB, queryParams []utils.FilterCondition, paginate pagination.Pagination) ([]map[string]interface{}, int, int, error)
	GetItemById(tx *gorm.DB, Id int) (masteritempayloads.ItemResponse, error)
	GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, error)
	GetItemCode(*gorm.DB, string) ([]map[string]interface{}, error)
	SaveItem(*gorm.DB, masteritempayloads.ItemResponse) (bool, error)
	ChangeStatusItem(tx *gorm.DB, Id int) (bool, error)
	AddItemDetail(*gorm.DB, int, masteritempayloads.ItemDetailRequest) *exceptionsss_test.BaseErrorResponse
	DeleteItemDetail(*gorm.DB, int, int) *exceptionsss_test.BaseErrorResponse
	GetItemDetailById(*gorm.DB, int, int) (masteritempayloads.ItemDetailRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
