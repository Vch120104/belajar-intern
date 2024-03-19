package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveMasterService interface {
	GetIncentiveMasterById(id int) masterpayloads.IncentiveMasterResponse
	SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) bool
	GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	ChangeStatusIncentiveMaster(Id int) bool
}
