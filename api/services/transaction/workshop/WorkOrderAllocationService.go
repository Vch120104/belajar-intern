package transactionworkshopservice

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
	"time"
)

type WorkOrderAllocationService interface {
	GetAll(companyCode int, foremanId int, date time.Time, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllocate(date time.Time, brandId int, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetWorkOrderAllocationHeaderData(companyCode string, foremanId int, techallocStartDate time.Time, vehicleBrandId int) (transactionworkshoppayloads.WorkOrderAllocationHeaderResult, *exceptions.BaseErrorResponse)
	GetAllocateDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveAllocateDetail(date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationDetailRequest, foremanId int, companyId int) (transactionworkshopentities.WorkOrderAllocationDetail, *exceptions.BaseErrorResponse)
	GetAssignTechnician(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	NewAssignTechnician(date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse)
	GetAssignTechnicianById(date time.Time, techId int, id int) (transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse, *exceptions.BaseErrorResponse)
	SaveAssignTechnician(date time.Time, techId int, id int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse)
}
