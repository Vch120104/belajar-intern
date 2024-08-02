package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LabourSellingPriceRepository interface {
	GetAllLabourSellingPrice(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetLabourSellingPriceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(tx *gorm.DB, headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}
