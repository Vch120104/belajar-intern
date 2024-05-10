package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountPercentService interface {
	GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusDiscountPercent(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
