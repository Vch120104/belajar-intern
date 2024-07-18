package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepository interface {
	GetAllOperationGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationGroupById(*gorm.DB, int) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	SaveOperationGroup(*gorm.DB, masteroperationpayloads.OperationGroupResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationGroup(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetOperationGroupByCode(*gorm.DB, string) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	GetAllOperationGroupIsActive(*gorm.DB) ([]masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
}
