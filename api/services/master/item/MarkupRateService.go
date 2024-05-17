package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupRateService interface {
	GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse)
	SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMarkupRate(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMarkupRateByMarkupMasterAndOrderType(MarkupMasterId int, OrderTypeId int) ([]masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse)
}
