package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LandedCostMasterService interface {
	GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdLandedCost(id int) (masteritempayloads.LandedCostMasterPayloads, *exceptions.BaseErrorResponse)
	SaveLandedCost(req masteritempayloads.LandedCostMasterPayloads) (bool, *exceptions.BaseErrorResponse)
	DeactivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse)
}
