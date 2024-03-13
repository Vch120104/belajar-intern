package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

)

type OperationKeyService interface {
	GetOperationKeyById(int) masteroperationpayloads.OperationkeyListResponse
	GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetOperationKeyName(masteroperationpayloads.OperationKeyRequest) masteroperationpayloads.OperationKeyNameResponse
	SaveOperationKey(masteroperationpayloads.OperationKeyResponse) bool
	ChangeStatusOperationKey(int) bool
}
