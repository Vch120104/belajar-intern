package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationCodeService interface {
	GetOperationCodeById(int) (masteroperationpayloads.OperationCodeResponse,*exceptionsss_test.BaseErrorResponse)
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination,*exceptionsss_test.BaseErrorResponse)
	SaveOperationCode(masteroperationpayloads.OperationCodeSave) (bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationCode(int) (bool,*exceptionsss_test.BaseErrorResponse)
}