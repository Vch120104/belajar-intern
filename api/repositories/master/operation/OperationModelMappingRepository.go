package masteroperationrepository

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationModelMappingRepository interface {
	GetOperationModelMappingById(*gorm.DB, int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse)
	GetOperationModelMappingByBrandModelOperationCode(*gorm.DB, masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse)
	GetOperationModelMappingLookup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveOperationModelMapping(*gorm.DB, masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusOperationModelMapping(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	SaveOperationModelMappingFrt(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	ActivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	SaveOperationModelMappingDocumentRequirement(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptions.BaseErrorResponse)
	DeactivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	ActivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	GetAllOperationDocumentRequirement(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationDocumentRequirementById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse)
	GetAllOperationFrt(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetOperationFrtById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse)
	SaveOperationLevel(tx *gorm.DB, request masteroperationpayloads.OperationLevelRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllOperationLevel(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOperationLevelById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationLevelByIdResponse, *exceptions.BaseErrorResponse)
	ActivateOperationLevel(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	DeactivateOperationLevel(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	DeleteOperationLevel(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse)
	// SaveOperationModelMappingAndFRT(tx *gorm.DB, requestHeader masteroperationpayloads.OperationModelMappingResponse, requestDetail masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse)
	GetOperationModelMappingLatestId(tx *gorm.DB) (int, *exceptions.BaseErrorResponse)
	UpdateOperationModelMapping(tx *gorm.DB, operationModelMappingId int, request masteroperationpayloads.OperationModelMappingUpdate) (masteroperationentities.OperationModelMapping, *exceptions.BaseErrorResponse)
}
