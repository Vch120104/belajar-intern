package transactionjpcbrepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OutstandingJobAllocationRepository interface {
	GetAllOutstandingJobAllocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByTypeIdOutstandingJobAllocation(tx *gorm.DB, referenceDocumentType string, referenceSystemNumber int) (transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse, *exceptions.BaseErrorResponse)
	SaveOutstandingJobAllocation(tx *gorm.DB, referenceDocumentType string, referenceSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationSaveRequest, operationCodeResponse masteroperationpayloads.OperationCodeResponse) (transactionworkshopentities.WorkOrderAllocation, transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest, *exceptions.BaseErrorResponse)
	UpdateOutstandingJobAllocation(tx *gorm.DB, techAllocSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest) (transactionjpcbpayloads.OutstandingJobAllocationUpdateResponse, *exceptions.BaseErrorResponse)
	ReCalculateTimeJob(tx *gorm.DB, techAllocSystemNumber int) *exceptions.BaseErrorResponse
}
