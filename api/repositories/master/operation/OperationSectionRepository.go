package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type OperationSectionRepository interface {
	GetOperationSectionById(*gorm.DB, int) (masteroperationpayloads.OperationSectionListResponse, error)
	GetSectionCodeByGroupId(*gorm.DB, int) ([]masteroperationpayloads.OperationSectionCodeResponse, error)
	GetOperationSectionName(*gorm.DB, int, string) (masteroperationpayloads.OperationSectionNameResponse, error)
	GetAllOperationSectionList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	SaveOperationSection(*gorm.DB, masteroperationpayloads.OperationSectionRequest) (bool, error)
	ChangeStatusOperationSection(*gorm.DB, int) (bool, error)
}
