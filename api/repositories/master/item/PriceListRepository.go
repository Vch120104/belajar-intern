package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"

	"gorm.io/gorm"
)

type PriceListRepository interface {
	GetPriceList(*gorm.DB, masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error)
	SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, error)
	GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListResponse, error)
	ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, error)
}
