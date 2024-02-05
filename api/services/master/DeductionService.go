package masterservice

import (
	"after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DeductionService interface {
	GetAllDeduction(filterCondition []utils.FilterCondition,pages pagination.Pagination)  pagination.Pagination
	GetByIdDeductionDetail(id int) masterpayloads.DeductionDetailResponse
	GetByIdDeductionList(id int, page int, limit int) payloads.ResponsePaginationHeader
	PostDeductionDetail(req masterpayloads.DeductionDetailResponse) masterpayloads.DeductionDetailResponse
	PostDeductionList(req masterpayloads.DeductionListResponse) masterpayloads.DeductionListResponse
}