package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationGroupService interface {
	GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	GetOperationGroupById(int) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	GetOperationGroupByCode(string) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	ChangeStatusOperationGroup(int) (bool, *exceptions.BaseErrorResponse)
	SaveOperationGroup(masteroperationpayloads.OperationGroupResponse) (bool, *exceptions.BaseErrorResponse)
}
