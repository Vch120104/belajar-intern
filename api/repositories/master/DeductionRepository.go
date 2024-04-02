package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionRepository interface {
	GetAllDeduction(*gorm.DB, []utils.FilterCondition,  pagination.Pagination) (pagination.Pagination, error)
	GetAllDeductionDetail(*gorm.DB,  pagination.Pagination,  int) (pagination.Pagination, error)
	GetDeductionById(*gorm.DB, int) (masterpayloads.DeductionListResponse, error)
	GetByIdDeductionDetail(*gorm.DB, int) (masterpayloads.DeductionDetailResponse, error)
	SaveDeductionList(*gorm.DB, masterpayloads.DeductionListResponse) (masterpayloads.DeductionListResponse, error)
	SaveDeductionDetail(*gorm.DB, masterpayloads.DeductionDetailResponse) (masterpayloads.DeductionDetailResponse, error)
	ChangeStatusDeduction(tx *gorm.DB, Id int) (bool, error)
}