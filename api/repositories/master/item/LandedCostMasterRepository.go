package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LandedCostMasterRepository interface {
	GetAllLandedCost(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdLandedCost(*gorm.DB, int) (masteritempayloads.LandedCostMasterPayloads, *exceptionsss_test.BaseErrorResponse)
	SaveLandedCost(*gorm.DB, masteritempayloads.LandedCostMasterPayloads) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateLandedCostmaster(*gorm.DB, string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateLandedCostMaster(*gorm.DB, string) (bool, *exceptionsss_test.BaseErrorResponse)
}
