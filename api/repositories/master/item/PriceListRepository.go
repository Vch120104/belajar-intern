package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"

	"gorm.io/gorm"
)

type PriceListRepository interface {
	GetPriceList(*gorm.DB, masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse)
	GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
