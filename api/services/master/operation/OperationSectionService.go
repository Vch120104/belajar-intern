package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	// "after-sales/api/utils"
)

type OperationSectionService interface {
	GetAllOperationSectionList([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationSectionById(int) (masteroperationpayloads.OperationSectionListResponse, *exceptions.BaseErrorResponse)
	GetSectionCodeByGroupId(int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptions.BaseErrorResponse)
	GetOperationSectionName(int, string) (masteroperationpayloads.OperationSectionNameResponse, *exceptions.BaseErrorResponse)
	SaveOperationSection(masteroperationpayloads.OperationSectionRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationSection(int) (bool, *exceptions.BaseErrorResponse)
}
