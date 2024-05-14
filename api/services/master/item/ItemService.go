package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemService interface {
	GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLookup(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse)
	GetItemById(Id int) (map[string]any, *exceptionsss_test.BaseErrorResponse)
	GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemCode(string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveItem(masteritempayloads.ItemRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItem(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptionsss_test.BaseErrorResponse)
	GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptionsss_test.BaseErrorResponse)
}
