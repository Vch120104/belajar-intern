package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
)

type ItemService interface {
	GetAllItem(filterCondition []utils.FilterCondition) ([]masteritempayloads.ItemLookup, *exceptionsss_test.BaseErrorResponse)
	GetAllItemLookup(map[string]string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse)
	GetItemCode(string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	SaveItem(masteritempayloads.ItemResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItem(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
