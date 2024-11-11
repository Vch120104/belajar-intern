package masteroperationrepository

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyRepository interface {
	GetOperationKeyById(*gorm.DB, int) (masteroperationpayloads.OperationkeyListResponse, *exceptions.BaseErrorResponse)
	GetOperationKeyName(*gorm.DB, masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptions.BaseErrorResponse)
	SaveOperationKey(*gorm.DB, masteroperationpayloads.OperationKeyResponse) (bool, *exceptions.BaseErrorResponse)
	GetAllOperationKeyList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusOperationKey(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	GetOperationKeyDropdown(tx *gorm.DB, operationGroupId int, operationSectionId int) ([]masteroperationpayloads.OperationKeyDropDown, *exceptions.BaseErrorResponse)
}
