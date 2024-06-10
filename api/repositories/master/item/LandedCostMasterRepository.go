package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LandedCostMasterRepository interface {
	GetAllLandedCost(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{},int,int,*exceptions.BaseErrorResponse)
	GetByIdLandedCost(*gorm.DB, int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveLandedCost(*gorm.DB, masteritempayloads.LandedCostMasterRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateLandedCostmaster(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	UpdateLandedCostMaster(tx *gorm.DB,id int, req masteritempayloads.LandedCostMasterUpdateRequest)(bool,*exceptions.BaseErrorResponse)
}
