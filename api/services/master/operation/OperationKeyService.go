package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationKeyService interface {
	GetOperationKeyById(int) (masteroperationpayloads.OperationkeyListResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationKeyName(masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationKey(masteroperationpayloads.OperationKeyResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationKey(int) (bool, *exceptionsss_test.BaseErrorResponse)
}
