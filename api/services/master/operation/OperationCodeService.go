package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationCodeService interface {
	GetOperationCodeById(int32) (masteroperationpayloads.OperationCodeResponse, error)
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
}
