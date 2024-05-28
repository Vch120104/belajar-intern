package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionRepository interface {
	GetAllDeduction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDeductionDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDeductionById(*gorm.DB, int) (masterpayloads.DeductionListResponse, *exceptions.BaseErrorResponse)
	GetByIdDeductionDetail(*gorm.DB, int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse)
	SaveDeductionList(*gorm.DB, masterpayloads.DeductionListResponse) (bool, *exceptions.BaseErrorResponse)
	SaveDeductionDetail(*gorm.DB, masterpayloads.DeductionDetailResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusDeduction(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
