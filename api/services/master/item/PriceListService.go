package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"

	"gorm.io/gorm"
)

type PriceListService interface {
	WithTrx(trxHandle *gorm.DB) PriceListService
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error)
	GetPriceListById(Id int) (masteritempayloads.PriceListResponse, error)
	SavePriceList(request masteritempayloads.PriceListResponse) (bool, error)
	ChangeStatusPriceList(Id int) (bool, error)
}
