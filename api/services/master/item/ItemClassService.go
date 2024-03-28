package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
)

type ItemClassService interface {
	GetAllItemClass(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, *exceptionsss_test.BaseErrorResponse)
	SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemClass(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
