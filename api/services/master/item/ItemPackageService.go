package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemPackageService interface {
	GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveItemPackage(request masteritempayloads.SaveItemPackageRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetItemPackageById(Id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemPackage(id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
