package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LandedCostMasterRepository interface {
	GetAllLandedCost(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetByIdLandedCost(*gorm.DB, int) (masteritempayloads.LandedCostMasterPayloads, error)
	SaveLandedCost(*gorm.DB, masteritempayloads.LandedCostMasterPayloads) (bool, error)
	DeactivateLandedCostmaster(*gorm.DB, string) (bool, error)
	ActivateLandedCostMaster(*gorm.DB, string) (bool, error)
}
