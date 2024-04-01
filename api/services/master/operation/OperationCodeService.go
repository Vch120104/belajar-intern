package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationCodeService interface {
	GetOperationCodeById(int) masteroperationpayloads.OperationCodeResponse
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination) pagination.Pagination
	SaveOperationCode(masteroperationpayloads.OperationCodeSave) bool
	ChangeStatusOperationCode(int) bool
}