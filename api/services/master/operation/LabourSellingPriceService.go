package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
)

type LabourSellingPriceService interface {
	GetLabourSellingPriceById(Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}