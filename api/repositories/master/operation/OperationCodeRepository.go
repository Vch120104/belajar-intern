package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepository interface {
	GetOperationCodeById(*gorm.DB,int) (masteroperationpayloads.OperationCodeResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationCode(*gorm.DB,[]utils.FilterCondition, pagination.Pagination)(pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SaveOperationCode(*gorm.DB,masteroperationpayloads.OperationCodeSave)(bool,*exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemSubstitute(*gorm.DB,int)(bool,*exceptionsss_test.BaseErrorResponse)
}
