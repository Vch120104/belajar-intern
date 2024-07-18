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
	SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse)
	GetPriceListById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	DeactivatePriceList (tx *gorm.DB, id string)(bool, *exceptions.BaseErrorResponse)
	ActivatePriceList (tx *gorm.DB, id string)(bool, *exceptions.BaseErrorResponse)
	DeletePriceList (tx*gorm.DB, id string)(bool,*exceptions.BaseErrorResponse)
	GetAllPriceListNew(tx *gorm.DB, filtercondition []utils.FilterCondition,pages pagination.Pagination)([]map[string]interface{},int,int,*exceptions.BaseErrorResponse)
}
