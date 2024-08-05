package masteroperationrepository

import (
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LabourSellingPriceRepository interface {
	GetLabourSellingPriceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(tx *gorm.DB, headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllSellingPrice(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptions.BaseErrorResponse)
	SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptions.BaseErrorResponse)
}
