package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveMasterRepository interface {
	GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, error)
	SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (bool, error)
	GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
	ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (bool, error)
}
