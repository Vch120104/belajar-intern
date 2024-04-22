package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepository interface {
	GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, error)
	SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (bool, error)
	GetItemPackageById(tx *gorm.DB, id int) ([]map[string]interface{}, error)
}
