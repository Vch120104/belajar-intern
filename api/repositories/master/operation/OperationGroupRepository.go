package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationGroupRepository
	GetAllOperationGroup([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetOperationGroupById(int) (masteroperationpayloads.OperationGroupResponse, error)
	SaveOperationGroup(masteroperationpayloads.OperationGroupResponse) (bool, error)
	ChangeStatusOperationGroup(int) (bool, error)
	GetOperationGroupByCode(string) (masteroperationpayloads.OperationGroupResponse, error)
	GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, error)
}
