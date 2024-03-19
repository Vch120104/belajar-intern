package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LandedCostMasterService interface {
	GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetByIdLandedCost(id int) masteritempayloads.LandedCostMasterPayloads
	SaveLandedCost(req masteritempayloads.LandedCostMasterPayloads) bool
	DeactivateLandedCostMaster(id string) bool
	ActivateLandedCostMaster(id string) bool
}