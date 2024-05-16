package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveMasterRepository interface {
	GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (masterentities.IncentiveMaster, *exceptionsss_test.BaseErrorResponse)
}
