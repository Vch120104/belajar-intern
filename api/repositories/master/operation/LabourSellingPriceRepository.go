package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type LabourSellingPriceRepository interface {
	GetLabourSellingPriceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(tx *gorm.DB, headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}