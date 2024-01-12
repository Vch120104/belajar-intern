package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyService interface {
	WithTrx(trxHandle *gorm.DB) OperationKeyService
	GetOperationKeyById(int) (masteroperationpayloads.OperationKeyResponse, error)
	GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetOperationKeyName(masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error)
	SaveOperationKey(masteroperationpayloads.OperationKeyResponse) (bool, error)
	ChangeStatusOperationKey(int) (bool, error)
}
