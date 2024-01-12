package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupService interface {
	WithTrx(trxHandle *gorm.DB) OperationGroupService
	GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, error)
	GetOperationGroupById(int) (masteroperationpayloads.OperationGroupResponse, error)
	GetOperationGroupByCode(string) (masteroperationpayloads.OperationGroupResponse, error)
	ChangeStatusOperationGroup(int) (bool, error)
	SaveOperationGroup(masteroperationpayloads.OperationGroupResponse) (bool, error)
}
