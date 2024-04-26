package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type OperationSectionRepository interface {
	GetOperationSectionById(*gorm.DB, int) (masteroperationpayloads.OperationSectionListResponse, *exceptionsss_test.BaseErrorResponse)
	GetSectionCodeByGroupId(*gorm.DB, int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationSectionName(*gorm.DB, int, string) (masteroperationpayloads.OperationSectionNameResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationSectionList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SaveOperationSection(*gorm.DB, masteroperationpayloads.OperationSectionRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationSection(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
