package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationCodeService interface {
	GetOperationCodeById(int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveOperationCode(masteroperationpayloads.OperationCodeSave) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationCode(int) (bool, *exceptions.BaseErrorResponse)
	GetOperationCodeByCode(string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
}
