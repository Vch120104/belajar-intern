package masteroperationservice

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationCodeService interface {
	GetOperationCodeById(int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllOperationCodeDropDown() ([]masteroperationpayloads.OperationCodeGetAll, *exceptions.BaseErrorResponse)
	SaveOperationCode(masteroperationpayloads.OperationCodeSave) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse)
	ChangeStatusOperationCode(int) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse)
	GetOperationCodeByCode(string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	UpdateItemCode(id int, req masteroperationpayloads.OperationCodeUpdate) (bool, *exceptions.BaseErrorResponse)
}
