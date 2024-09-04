package transactionjpcbservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"
)

type JobAllocationService interface {
	GetAllJobAllocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetJobAllocationById(technicianAllocationSystemNumber int) (transactionjpcbpayloads.GetJobAllocationByIdResponse, *exceptions.BaseErrorResponse)
	UpdateJobAllocation(technicianAllocationSystemNumber int, req transactionjpcbpayloads.JobAllocationUpdateRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse)
	DeleteJobAllocation(technicianAllocationSystemNumber int) (bool, *exceptions.BaseErrorResponse)
}
