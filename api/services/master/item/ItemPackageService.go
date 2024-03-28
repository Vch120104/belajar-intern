package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemPackageService interface {
	GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	SaveItemPackage(request masteritempayloads.SaveItemPackageRequest) bool
	GetItemPackageById(Id int) []map[string]interface{}
}
