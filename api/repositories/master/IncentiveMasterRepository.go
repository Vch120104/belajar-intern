package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveMasterRepository interface {
	GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
	UpdateIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest, Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
	GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
}
