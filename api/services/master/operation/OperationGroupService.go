package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationGroupService interface {
	GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetAllOperationGroupIsActive() []masteroperationpayloads.OperationGroupResponse
	GetOperationGroupById(int) masteroperationpayloads.OperationGroupResponse
	GetOperationGroupByCode(string) masteroperationpayloads.OperationGroupResponse
	ChangeStatusOperationGroup(int) bool
	SaveOperationGroup(masteroperationpayloads.OperationGroupResponse) bool
}
