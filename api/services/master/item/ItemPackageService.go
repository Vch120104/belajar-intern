package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemPackageService interface {
	GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveItemPackage(request masteritempayloads.SaveItemPackageRequest) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageById(Id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(id int) (bool, *exceptions.BaseErrorResponse)
}
