package masteroperationrepository

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationModelMappingRepository interface {
	WithTrx(trxHandle *gorm.DB) OperationModelMappingRepository
	GetOperationModelMappingById(int) (masteroperationpayloads.OperationModelMappingResponse, error)
	GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, error)
	GetOperationModelMappingLookup([]utils.FilterCondition) ([]map[string]interface{}, error)
	SaveOperationModelMapping(masteroperationpayloads.OperationModelMappingResponse) (bool, error)
	ChangeStatusOperationModelMapping(int) (bool, error)
}
