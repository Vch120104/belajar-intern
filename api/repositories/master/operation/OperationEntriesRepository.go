package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationEntriesRepository interface {
	GetAllOperationEntries(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationEntriesById(*gorm.DB, int) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationEntriesName(*gorm.DB, masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationEntries(*gorm.DB, masteroperationpayloads.OperationEntriesResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationEntries(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}