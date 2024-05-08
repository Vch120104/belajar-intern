package masteroperationservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationModelMappingService interface {
	GetOperationModelMappingById(int) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationModelMappingLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMapping(masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationModelMapping(int) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMappingFrt(request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateOperationFrt(id string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateOperationFrt(id string) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMappingDocumentRequirement(request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateOperationDocumentRequirement(id string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateOperationDocumentRequirement(id string) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationDocumentRequirement(id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationDocumentRequirementById(id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationFrt(id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationFrtById(id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptionsss_test.BaseErrorResponse)
}
