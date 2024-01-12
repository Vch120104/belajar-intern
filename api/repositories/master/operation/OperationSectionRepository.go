package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type OperationSectionRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationSectionRepository
	GetAllOperationSection() ([]masteroperationpayloads.OperationSectionResponse, error)
	GetOperationSectionById(int) (masteroperationpayloads.OperationSectionResponse, error)
	GetSectionCodeByGroupId(string) ([]masteroperationpayloads.OperationSectionCodeResponse, error)
	GetOperationSectionName(int, string) (masteroperationpayloads.OperationSectionNameResponse, error)
	GetAllOperationSectionList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	SaveOperationSection(masteroperationpayloads.OperationSectionRequest) (bool, error)
	ChangeStatusOperationSection(int) (bool, error)
}
