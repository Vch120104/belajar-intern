package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionRepository interface {
	GetAllDeduction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllDeductionDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetDeductionById(*gorm.DB, int, pagination.Pagination) (masterpayloads.DeductionById, *exceptions.BaseErrorResponse)
	GetByIdDeductionDetail(*gorm.DB, int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse)
	SaveDeductionList(*gorm.DB, masterpayloads.DeductionListResponse) (masterentities.DeductionList, *exceptions.BaseErrorResponse)
	SaveDeductionDetail(*gorm.DB, masterpayloads.DeductionDetailResponse, int) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse)
	ChangeStatusDeduction(*gorm.DB,  int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	UpdateDeductionDetail( *gorm.DB, int, masterpayloads.DeductionDetailUpdate)(masterentities.DeductionDetail,*exceptions.BaseErrorResponse)
}
