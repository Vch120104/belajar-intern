package masteroperationservice

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationModelMappingService interface {
	WithTrx(trxHandle *gorm.DB) OperationModelMappingService
	GetOperationModelMappingById(int) (masteroperationpayloads.OperationModelMappingResponse, error)
	GetOperationModelMappingLookup(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, error)
	SaveOperationModelMapping(masteroperationpayloads.OperationModelMappingResponse) (bool, error)
	ChangeStatusOperationModelMapping(int) (bool, error)
}
