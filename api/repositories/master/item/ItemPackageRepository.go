package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepository interface {
	GetAllItemPackage(tx *gorm.DB, internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveItemPackage(tx *gorm.DB, request masteritempayloads.SaveItemPackageRequest) (masteritementities.ItemPackage, *exceptions.BaseErrorResponse)
	GetItemPackageById(tx *gorm.DB, id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse)
	GetItemPackageByItemId(tx *gorm.DB, itemId int) (masteritempayloads.GetItemPackageItemResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackage(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageByCode(tx *gorm.DB, itemPackageCode string) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse)
}
