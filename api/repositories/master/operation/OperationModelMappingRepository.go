package masteroperationrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationModelMappingRepository interface {
	GetOperationModelMappingById(*gorm.DB, int) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationModelMappingByBrandModelOperationCode(*gorm.DB, masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse)
	GetOperationModelMappingLookup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMapping(*gorm.DB, masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusOperationModelMapping(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMappingFrt(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveOperationModelMappingDocumentRequirement(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationDocumentRequirement(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationDocumentRequirementById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllOperationFrt(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetOperationFrtById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptionsss_test.BaseErrorResponse)
}
