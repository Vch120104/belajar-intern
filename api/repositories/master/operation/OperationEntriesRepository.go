package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationEntriesRepository interface {
	GetAllOperationEntries(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationEntriesById(*gorm.DB, int) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse)
	GetOperationEntriesName(*gorm.DB, masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse)
	SaveOperationEntries(*gorm.DB, masteroperationpayloads.OperationEntriesResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationEntries(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
