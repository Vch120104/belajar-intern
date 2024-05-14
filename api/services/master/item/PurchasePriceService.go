package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PurchasePriceService interface {
	DeletePurchasePrice(id int) *exceptionsss_test.BaseErrorResponse
	GetPurchasePriceById(id int) (masteritempayloads.PurchasePriceRequest, *exceptionsss_test.BaseErrorResponse)
	SavePurchasePrice(masteritempayloads.PurchasePriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	AddPurchasePrice(masteritempayloads.PurchasePriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusPurchasePrice(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
