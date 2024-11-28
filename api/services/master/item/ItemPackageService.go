package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemPackageService interface {
	GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveItemPackage(request masteritempayloads.SaveItemPackageRequest) (masteritementities.ItemPackage, *exceptions.BaseErrorResponse)
	GetItemPackageById(Id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(id int) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageByCode(itemPackageCode string) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse)
}
