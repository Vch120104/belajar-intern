package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveGroupService interface {
	GetAllIncentiveGroup([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse)
	GetIncentiveGroupById(int) (masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveGroup(masterpayloads.IncentiveGroupResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusIncentiveGroup(int) (bool, *exceptions.BaseErrorResponse)
	UpdateIncentiveGroup(req masterpayloads.UpdateIncentiveGroupRequest) (bool, *exceptions.BaseErrorResponse)
}
