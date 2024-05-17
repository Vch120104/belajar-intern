package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveGroupService interface {
	GetAllIncentiveGroup([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetIncentiveGroupById(int) (masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveGroup(masterpayloads.IncentiveGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusIncentiveGroup(int) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateIncentiveGroup(req masterpayloads.UpdateIncentiveGroupRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllIncentiveGroupDropDown() ([]masterpayloads.IncentiveGroupDropDown, *exceptionsss_test.BaseErrorResponse)
}
