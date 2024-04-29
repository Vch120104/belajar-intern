package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	// "after-sales/api/utils"
)

type OperationSectionService interface {
	GetAllOperationSectionList([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationSectionById(int) (masteroperationpayloads.OperationSectionListResponse, *exceptionsss_test.BaseErrorResponse)
	GetSectionCodeByGroupId(int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationSectionName(int, string) (masteroperationpayloads.OperationSectionNameResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationSection(masteroperationpayloads.OperationSectionRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationSection(int) (bool, *exceptionsss_test.BaseErrorResponse)
}
