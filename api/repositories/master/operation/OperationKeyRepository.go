package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationKeyRepository
	GetOperationKeyById(int) (masteroperationpayloads.OperationKeyResponse, error)
	GetOperationKeyName(masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error)
	SaveOperationKey(masteroperationpayloads.OperationKeyResponse) (bool, error)
	GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatusOperationKey(int) (bool, error)
}
