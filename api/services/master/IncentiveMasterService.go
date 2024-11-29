package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveMasterService interface {
	GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
	UpdateIncentiveMaster(req masterpayloads.IncentiveMasterRequest, id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
	GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusIncentiveMaster(Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse)
}
