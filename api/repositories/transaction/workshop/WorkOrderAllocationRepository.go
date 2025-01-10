package transactionworkshoprepository

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
	"time"

	"gorm.io/gorm"
)

type WorkOrderAllocationRepository interface {
	GetAll(tx *gorm.DB, companyCode int, foremanId int, date time.Time, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetWorkOrderAllocationHeaderData(tx *gorm.DB, companyId int, foremanId int, techallocStartDate time.Time) (transactionworkshoppayloads.WorkOrderAllocationHeaderResult, *exceptions.BaseErrorResponse)
	GetAllocate(tx *gorm.DB, companyId int, date time.Time, foremanId int, brandId int, workOrderSystemNumber int, pages pagination.Pagination) (transactionworkshoppayloads.WorkOrderAllocationResponse, *exceptions.BaseErrorResponse)
	WorkOrderAllocationGR(tx *gorm.DB, companyId int, date time.Time, foremanId int, brandId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllocateDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveAllocateDetail(tx *gorm.DB, date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationDetailRequest, foremanId int, companyId int) (transactionworkshopentities.WorkOrderAllocationDetail, *exceptions.BaseErrorResponse)
	GetAssignTechnician(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	NewAssignTechnician(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse)
	GetAssignTechnicianById(tx *gorm.DB, date time.Time, techId int, id int) (transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse, *exceptions.BaseErrorResponse)
	SaveAssignTechnician(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse)
}
