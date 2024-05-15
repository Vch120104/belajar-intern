package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveMasterService interface {
	GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusIncentiveMaster(Id int) (masterentities.IncentiveMaster, *exceptionsss_test.BaseErrorResponse)
}
