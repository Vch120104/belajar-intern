package masteroperationservice

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type OperationModelMappingService interface {
	GetOperationModelMappingById(int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse)
	GetOperationModelMappingLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse)
	SaveOperationModelMapping(masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationModelMapping(int) (bool, *exceptions.BaseErrorResponse)
	SaveOperationModelMappingFrt(request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse)
	SaveOperationModelMappingDocumentRequirement(request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse)
	GetAllOperationDocumentRequirement(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationDocumentRequirementById(id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse)
	GetAllOperationFrt(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationFrtById(id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse)
}
