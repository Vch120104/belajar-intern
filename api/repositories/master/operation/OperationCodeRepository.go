package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepository interface {
	GetOperationCodeById(*gorm.DB, int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	GetAllOperationCode(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveOperationCode(*gorm.DB, masteroperationpayloads.OperationCodeSave) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemSubstitute(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetOperationCodeByCode(*gorm.DB, string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
}
