package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
)

type PriceListService interface {
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptionsss_test.BaseErrorResponse)
	GetPriceListById(Id int) (masteritempayloads.PriceListResponse, *exceptionsss_test.BaseErrorResponse)
	SavePriceList(request masteritempayloads.PriceListResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusPriceList(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
