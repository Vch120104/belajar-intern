package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyRepository interface {
	GetOperationKeyById(*gorm.DB, int) (masteroperationpayloads.OperationkeyListResponse, error)
	GetOperationKeyName(*gorm.DB, masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error)
	SaveOperationKey(*gorm.DB, masteroperationpayloads.OperationKeyResponse) (bool, error)
	GetAllOperationKeyList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	ChangeStatusOperationKey(*gorm.DB, int) (bool, error)
}
