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
	SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
}
