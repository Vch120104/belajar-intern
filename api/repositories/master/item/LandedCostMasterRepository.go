package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LandedCostMasterRepository interface {
	GetAllLandedCost(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdLandedCost(tx *gorm.DB, id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveLandedCost(tx *gorm.DB, req masteritempayloads.LandedCostMasterRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse)
	DeactivateLandedCostmaster(tx *gorm.DB, id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(tx *gorm.DB, id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	UpdateLandedCostMaster(tx *gorm.DB, id int, req masteritempayloads.LandedCostMasterUpdateRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse)
}
