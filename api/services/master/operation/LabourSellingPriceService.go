package masteroperationservice

import (
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LabourSellingPriceService interface {
	GetLabourSellingPriceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllSellingPrice(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptions.BaseErrorResponse)
	SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptions.BaseErrorResponse)
}
