package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationEntriesService interface {
	GetAllOperationEntries([]utils.FilterCondition, pagination.Pagination) pagination.Pagination
	GetOperationEntriesById(int) masteroperationpayloads.OperationEntriesResponse
	GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest) masteroperationpayloads.OperationEntriesResponse
	SaveOperationEntries(masteroperationpayloads.OperationEntriesResponse) bool
	ChangeStatusOperationEntries(Id int) bool
}
