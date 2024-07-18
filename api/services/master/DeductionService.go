package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DeductionService interface {
	GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdDeductionDetail(id int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse)
	GetDeductionById(Id int) (masterpayloads.DeductionListResponse, *exceptions.BaseErrorResponse)
	GetAllDeductionDetail(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	PostDeductionDetail(req masterpayloads.DeductionDetailResponse) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse)
	PostDeductionList(req masterpayloads.DeductionListResponse) (masterentities	.DeductionList, *exceptions.BaseErrorResponse)
	ChangeStatusDeduction(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	UpdateDeductionDetail(id int, req masterpayloads.DeductionDetailUpdate)(masterentities.DeductionDetail,*exceptions.BaseErrorResponse)
}
