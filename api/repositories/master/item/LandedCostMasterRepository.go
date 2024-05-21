package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LandedCostMasterRepository interface {
	GetAllLandedCost(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdLandedCost(*gorm.DB, int) (masteritempayloads.LandedCostMasterPayloads, *exceptions.BaseErrorResponse)
	SaveLandedCost(*gorm.DB, masteritempayloads.LandedCostMasterPayloads) (bool, *exceptions.BaseErrorResponse)
	DeactivateLandedCostmaster(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
}
