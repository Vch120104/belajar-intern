package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationCodeRepository
	GetOperationCodeById(int32) (masteroperationpayloads.OperationCodeResponse, error)
	GetAllOperationCode([]utils.FilterCondition, pagination.Pagination)(pagination.Pagination, error)
}
