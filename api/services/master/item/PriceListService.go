package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
)

type PriceListService interface {
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) []masteritempayloads.PriceListResponse
	GetPriceListById(Id int) masteritempayloads.PriceListResponse
	SavePriceList(request masteritempayloads.PriceListResponse) bool
	ChangeStatusPriceList(Id int) bool
}
