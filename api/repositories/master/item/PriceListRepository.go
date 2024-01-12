package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"

	"gorm.io/gorm"
)

type PriceListRepository interface {
	WithTrx(trxHandle *gorm.DB) PriceListRepository
	GetPriceList(masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error)
	SavePriceList(request masteritempayloads.PriceListResponse) (bool, error)
	GetPriceListById(Id int) (masteritempayloads.PriceListResponse, error)
	ChangeStatusPriceList(Id int) (bool, error)
}
