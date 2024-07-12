package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LandedCostMasterService interface {
	GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse)
	GetByIdLandedCost(id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveLandedCost(req masteritempayloads.LandedCostMasterRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse)
	DeactivateLandedCostMaster(id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	UpdateLandedCostMaster(id int,req masteritempayloads.LandedCostMasterUpdateRequest)(masteritementities.LandedCost,*exceptions.BaseErrorResponse)
}
