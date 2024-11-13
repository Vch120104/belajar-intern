package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepository interface {
	GetAllOperationGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationGroupById(tx *gorm.DB, id int) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	SaveOperationGroup(tx *gorm.DB, req masteroperationpayloads.OperationGroupResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationGroup(tx *gorm.DB, oprId int) (bool, *exceptions.BaseErrorResponse)
	GetOperationGroupByCode(tx *gorm.DB, Code string) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse)
	GetOperationGroupDropDown(tx *gorm.DB) ([]masteroperationpayloads.OperationGroupDropDownResponse, *exceptions.BaseErrorResponse)
}
