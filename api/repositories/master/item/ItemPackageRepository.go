package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepository interface {
	GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse)
	SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageById(tx *gorm.DB, id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
}
