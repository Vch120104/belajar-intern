package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LandedCostMasterService interface {
	GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse)
	GetByIdLandedCost(id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	SaveLandedCost(req masteritempayloads.LandedCostMasterRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse)
	UpdateLandedCostMaster(id int,req masteritempayloads.LandedCostMasterUpdateRequest)(bool,*exceptions.BaseErrorResponse)
}
