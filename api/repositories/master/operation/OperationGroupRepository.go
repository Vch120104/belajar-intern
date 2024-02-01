package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepository interface {
	GetAllOperationGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetOperationGroupById(*gorm.DB, int) (masteroperationpayloads.OperationGroupResponse, error)
	SaveOperationGroup(*gorm.DB, masteroperationpayloads.OperationGroupResponse) (bool, error)
	ChangeStatusOperationGroup(*gorm.DB, int) (bool, error)
	GetOperationGroupByCode(*gorm.DB, string) (masteroperationpayloads.OperationGroupResponse, error)
	GetAllOperationGroupIsActive(*gorm.DB) ([]masteroperationpayloads.OperationGroupResponse, error)
}
