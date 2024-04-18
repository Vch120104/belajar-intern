package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
)

type LabourSellingPriceService interface {
	SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}