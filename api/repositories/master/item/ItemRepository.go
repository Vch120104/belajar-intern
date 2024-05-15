package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLookup(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	GetItemById(tx *gorm.DB, Id int) (map[string]any, *exceptionsss_test.BaseErrorResponse)
	GetItemWithMultiId(tx *gorm.DB, MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemCode(*gorm.DB, string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveItem(*gorm.DB, masteritempayloads.ItemRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItem(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetUomTypeDropDown(tx *gorm.DB) ([]masteritempayloads.UomTypeDropdownResponse, *exceptionsss_test.BaseErrorResponse)
	GetUomDropDown(tx *gorm.DB, uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptionsss_test.BaseErrorResponse)
	AddItemDetail(*gorm.DB, int, masteritempayloads.ItemDetailRequest) *exceptionsss_test.BaseErrorResponse
	DeleteItemDetail(*gorm.DB, int, int) *exceptionsss_test.BaseErrorResponse
	GetItemDetailById(*gorm.DB, int, int) (masteritempayloads.ItemDetailRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllItemDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
