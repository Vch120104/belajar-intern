package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
)

type PriceListService interface {
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	GetPriceListById(Id int) (masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	SavePriceList(request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(Id int) (bool, *exceptions.BaseErrorResponse)
}
