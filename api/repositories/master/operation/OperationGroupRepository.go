package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepository interface {
	GetAllOperationGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationGroupById(*gorm.DB, int) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationGroup(*gorm.DB, masteroperationpayloads.OperationGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationGroup(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetOperationGroupByCode(*gorm.DB, string) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationGroupIsActive(*gorm.DB) ([]masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse)
}
