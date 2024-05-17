package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LandedCostMasterService interface {
	GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	GetByIdLandedCost(id int) (masteritempayloads.LandedCostMasterPayloads,*exceptionsss_test.BaseErrorResponse)
	SaveLandedCost(req masteritempayloads.LandedCostMasterPayloads) (bool,*exceptionsss_test.BaseErrorResponse)
	DeactivateLandedCostMaster(id string) (bool,*exceptionsss_test.BaseErrorResponse)
	ActivateLandedCostMaster(id string) (bool,*exceptionsss_test.BaseErrorResponse)
}
