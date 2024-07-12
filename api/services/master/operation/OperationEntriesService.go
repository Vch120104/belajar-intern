package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationEntriesService interface {
	GetAllOperationEntries([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationEntriesById(int) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse)
	GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse)
	SaveOperationEntries(masteroperationpayloads.OperationEntriesResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationEntries(Id int) (bool, *exceptions.BaseErrorResponse)
}
