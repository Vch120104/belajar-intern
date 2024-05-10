package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationGroupService interface {
	GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationGroupById(int) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationGroupByCode(string) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationGroup(int) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveOperationGroup(masteroperationpayloads.OperationGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse)
}
