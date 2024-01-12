package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"

	"gorm.io/gorm"
)

type OperationEntriesRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationEntriesRepository
	GetOperationEntriesById(int32) (masteroperationpayloads.OperationEntriesResponse, error)
	GetOperationEntriesName(request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error)
	SaveOperationEntries(req masteroperationpayloads.OperationEntriesResponse) (bool, error)
	ChangeStatusOperationEntries(int) (bool, error)
}
