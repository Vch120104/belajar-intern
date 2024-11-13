package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationKeyService interface {
	GetOperationKeyById(int) (masteroperationpayloads.OperationkeyListResponse, *exceptions.BaseErrorResponse)
	GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationKeyName(masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptions.BaseErrorResponse)
	SaveOperationKey(masteroperationpayloads.OperationKeyResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationKey(int) (bool, *exceptions.BaseErrorResponse)
	GetOperationKeyDropdown(operationGroupId int, operationSectionId int) ([]masteroperationpayloads.OperationKeyDropDown, *exceptions.BaseErrorResponse)
}
