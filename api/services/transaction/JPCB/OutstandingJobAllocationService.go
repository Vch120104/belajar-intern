package transactionjpcbservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type OutstandingJobAllocationService interface {
	GetAllOutstandingJobAllocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByTypeIdOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int) (transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse, *exceptions.BaseErrorResponse)
	SaveOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationSaveRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse)
}
