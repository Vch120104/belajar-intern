package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type OperationSectionRepository interface {
	GetOperationSectionById(*gorm.DB, int) (masteroperationpayloads.OperationSectionListResponse, *exceptions.BaseErrorResponse)
	GetSectionCodeByGroupId(*gorm.DB, int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptions.BaseErrorResponse)
	GetOperationSectionName(*gorm.DB, int, string) (masteroperationpayloads.OperationSectionNameResponse, *exceptions.BaseErrorResponse)
	GetAllOperationSectionList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveOperationSection(*gorm.DB, masteroperationpayloads.OperationSectionRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationSection(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetOperationSectionDropDown(tx *gorm.DB, operationGroupId int) ([]masteroperationpayloads.OperationSectionDropDown, *exceptions.BaseErrorResponse)
}
