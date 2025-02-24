package masteroperationservice

import (
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LabourSellingPriceService interface {
	GetLabourSellingPriceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllSellingPriceDetailByHeaderId(headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSellingPriceDetailById(detailId int) (masteroperationpayloads.LabourSellingPriceDetailbyIdResponse, *exceptions.BaseErrorResponse)
	GetAllSellingPrice(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (int, *exceptions.BaseErrorResponse)
	SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (int, *exceptions.BaseErrorResponse)
	Duplicate(headerId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveDuplicate(req masteroperationpayloads.SaveDuplicateLabourSellingPrice) (masteroperationpayloads.LabourSellingSaveDuplicateResp, *exceptions.BaseErrorResponse)
	DeleteLabourSellingPriceDetail(iddet []int) (bool, *exceptions.BaseErrorResponse)
}
