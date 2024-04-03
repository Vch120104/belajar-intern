package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DeductionService interface {
	GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdDeductionDetail(id int) (masterpayloads.DeductionDetailResponse, *exceptionsss_test.BaseErrorResponse)
	GetDeductionById(Id int) (masterpayloads.DeductionListResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllDeductionDetail(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	PostDeductionDetail(req masterpayloads.DeductionDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	PostDeductionList(req masterpayloads.DeductionListResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusDeduction(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
