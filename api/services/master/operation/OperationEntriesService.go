package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"

	"gorm.io/gorm"
)

type OperationEntriesService interface {
	WithTrx(trxHandle *gorm.DB) OperationEntriesService
	GetOperationEntriesById(int32) (masteroperationpayloads.OperationEntriesResponse, error)
	GetOperationEntriesName(masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error)
	SaveOperationEntries(masteroperationpayloads.OperationEntriesResponse) (bool, error)
	ChangeStatusOperationEntries(Id int) (bool, error)
}
