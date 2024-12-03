package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PriceListRepository interface {
	GetPriceList(*gorm.DB, masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	SavePriceList(tx *gorm.DB, request masteritempayloads.SavePriceListMultiple) (int, *exceptions.BaseErrorResponse)
	GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListGetbyId, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	DeactivatePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	ActivatePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	DeletePriceList(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	// CheckPriceListAlreadyExist FOR UPLOAD TEMPLATE
	CheckPriceListExist(tx *gorm.DB, itemId int, brandId int, currencyId int, date string, companyId int) (bool, *exceptions.BaseErrorResponse)
	CheckPriceListItem(tx *gorm.DB, itemGroupId int, brandId int, currencyId int, date string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllPriceListNew(tx *gorm.DB, filtercondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	Duplicate(tx *gorm.DB, itemGroupId int, brandId int, currencyId int, date string) ([]masteritempayloads.PriceListItemResponses, *exceptions.BaseErrorResponse)
}
