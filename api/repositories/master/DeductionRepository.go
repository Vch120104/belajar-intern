package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionRepository interface {
	GetAllDeduction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllDeductionDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetDeductionById(*gorm.DB, int) (masterpayloads.DeductionListResponse, *exceptionsss_test.BaseErrorResponse)
	GetByIdDeductionDetail(*gorm.DB, int) (masterpayloads.DeductionDetailResponse, *exceptionsss_test.BaseErrorResponse)
	SaveDeductionList(*gorm.DB, masterpayloads.DeductionListResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveDeductionDetail(*gorm.DB, masterpayloads.DeductionDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusDeduction(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
