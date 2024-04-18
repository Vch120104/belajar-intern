package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"

	"gorm.io/gorm"
)

type LabourSellingPriceRepository interface {
	SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}