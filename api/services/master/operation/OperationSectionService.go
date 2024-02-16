package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	// "after-sales/api/utils"
)

type OperationSectionService interface {
	GetAllOperationSectionList([]utils.FilterCondition, pagination.Pagination) pagination.Pagination
	GetOperationSectionById(int) masteroperationpayloads.OperationSectionListResponse
	GetSectionCodeByGroupId(int) []masteroperationpayloads.OperationSectionCodeResponse
	GetOperationSectionName(int, string) masteroperationpayloads.OperationSectionNameResponse
	SaveOperationSection(masteroperationpayloads.OperationSectionRequest) bool
	ChangeStatusOperationSection(int) bool
}
