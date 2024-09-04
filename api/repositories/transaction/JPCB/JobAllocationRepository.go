package transactionjpcbrepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type JobAllocationRepository interface {
	GetAllJobAllocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetJobAllocationById(tx *gorm.DB, technicianAllocationSystemNumber int) (transactionjpcbpayloads.GetJobAllocationByIdResponse, *exceptions.BaseErrorResponse)
	UpdateJobAllocation(tx *gorm.DB, technicianAllocationSystemNumber int, req transactionjpcbpayloads.JobAllocationUpdateRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse)
	DeleteJobAllocation(tx *gorm.DB, technicianAllocationSystemNumber int) (bool, *exceptions.BaseErrorResponse)
}
