package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/utils"
	"after-sales/api/payloads/pagination"


	"gorm.io/gorm"
)

type OperationEntriesRepository interface {
	GetAllOperationEntries(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetOperationEntriesById(*gorm.DB, int) (masteroperationpayloads.OperationEntriesResponse, error)
	GetOperationEntriesName(*gorm.DB, masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error)
	SaveOperationEntries(*gorm.DB, masteroperationpayloads.OperationEntriesResponse) (bool, error)
	ChangeStatusOperationEntries(*gorm.DB, int) (bool, error)
}
