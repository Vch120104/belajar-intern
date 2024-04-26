package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyRepository interface {
	GetOperationKeyById(*gorm.DB, int) (masteroperationpayloads.OperationkeyListResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationKeyName(*gorm.DB, masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationKey(*gorm.DB, masteroperationpayloads.OperationKeyResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationKeyList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationKey(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
