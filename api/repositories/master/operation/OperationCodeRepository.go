package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepository interface {
	GetOperationCodeById(*gorm.DB,int) (masteroperationpayloads.OperationCodeResponse, error)
	GetAllOperationCode(*gorm.DB,[]utils.FilterCondition, pagination.Pagination)(pagination.Pagination, error)
	SaveOperationCode(*gorm.DB,masteroperationpayloads.OperationCodeSave)(bool,error)
	ChangeStatusItemSubstitute(*gorm.DB,int)(bool,error)
}
