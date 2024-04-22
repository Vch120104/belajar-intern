package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationEntriesService interface {
	GetAllOperationEntries([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationEntriesById(int) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationEntries(masteroperationpayloads.OperationEntriesResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationEntries(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}