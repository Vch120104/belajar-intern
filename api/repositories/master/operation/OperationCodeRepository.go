package masteroperationrepository

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepository interface {
	GetOperationCodeById(*gorm.DB, int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	GetAllOperationCode(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveOperationCode(*gorm.DB, masteroperationpayloads.OperationCodeSave) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse)
	ChangeStatusItemCode(*gorm.DB, int) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse)
	GetOperationCodeByCode(*gorm.DB, string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse)
	UpdateItemCode(tx *gorm.DB, id int, req masteroperationpayloads.OperationCodeUpdate)(masteroperationentities.OperationCode,*exceptions.BaseErrorResponse)
}
