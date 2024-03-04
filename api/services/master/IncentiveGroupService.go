package masterservice

import (
	masterpayloads "after-sales/api/payloads/master" 
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveGroupService interface {
	GetAllIncentiveGroup([]utils.FilterCondition,  pagination.Pagination) pagination.Pagination
	GetAllIncentiveGroupIsActive() []masterpayloads.IncentiveGroupResponse
	GetIncentiveGroupById(int) masterpayloads.IncentiveGroupResponse
	SaveIncentiveGroup(masterpayloads.IncentiveGroupResponse) bool
	ChangeStatusIncentiveGroup(int) bool
}
