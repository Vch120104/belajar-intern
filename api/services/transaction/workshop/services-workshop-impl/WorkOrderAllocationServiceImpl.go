package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"time"

	"after-sales/api/utils"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkOrderAllocationServiceImpl struct {
	WorkOrderAllocationRepository transactionworkshoprepository.WorkOrderAllocationRepository
	DB                            *gorm.DB
	RedisClient                   *redis.Client // Redis client
}

func OpenWorkOrderAllocationServiceImpl(WorkOrderAllocationRepo transactionworkshoprepository.WorkOrderAllocationRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.WorkOrderAllocationService {
	return &WorkOrderAllocationServiceImpl{
		WorkOrderAllocationRepository: WorkOrderAllocationRepo,
		DB:                            db,
		RedisClient:                   redisClient,
	}
}

func (s *WorkOrderAllocationServiceImpl) GetAll(companyCode int, foremanId int, date time.Time, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.WorkOrderAllocationRepository.GetAll(tx, companyCode, foremanId, date, filterCondition)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *WorkOrderAllocationServiceImpl) GetWorkOrderAllocationHeaderData(companyCode string, foremanId int, techallocStartDate time.Time, vehicleBrandId int) (transactionworkshoppayloads.WorkOrderAllocationHeaderResult, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.WorkOrderAllocationRepository.GetWorkOrderAllocationHeaderData(tx, companyCode, foremanId, techallocStartDate, vehicleBrandId)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *WorkOrderAllocationServiceImpl) GetAllocate(date time.Time, brandId int, woSysNum int) (transactionworkshoppayloads.WorkOrderAllocationResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.WorkOrderAllocationRepository.GetAllocate(tx, date, brandId, woSysNum)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *WorkOrderAllocationServiceImpl) GetAllocateDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.WorkOrderAllocationRepository.GetAllocateDetail(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderAllocationServiceImpl) GetAssignTechnician(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.WorkOrderAllocationRepository.GetAssignTechnician(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderAllocationServiceImpl) NewAssignTechnician(date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	entity, err := s.WorkOrderAllocationRepository.NewAssignTechnician(tx, date, techId, request)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	utils.RefreshCaches(ctx, "assign_technician")

	return entity, nil
}

func (s *WorkOrderAllocationServiceImpl) GetAssignTechnicianById(date time.Time, techId int, id int) (transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.WorkOrderAllocationRepository.GetAssignTechnicianById(tx, date, techId, id)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *WorkOrderAllocationServiceImpl) SaveAssignTechnician(date time.Time, techId int, id int, request transactionworkshoppayloads.WorkOrderAllocationAssignTechnicianRequest) (transactionworkshopentities.AssignTechnician, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	entity, err := s.WorkOrderAllocationRepository.SaveAssignTechnician(tx, date, techId, id, request)
	if err != nil {
		return transactionworkshopentities.AssignTechnician{}, err
	}

	utils.RefreshCaches(ctx, "assign_technician")

	return entity, nil
}

func (s *WorkOrderAllocationServiceImpl) SaveAllocateDetail(date time.Time, techId int, request transactionworkshoppayloads.WorkOrderAllocationDetailRequest, foremanId int, companyId int) (transactionworkshopentities.WorkOrderAllocationDetail, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	entity, err := s.WorkOrderAllocationRepository.SaveAllocateDetail(tx, date, techId, request, foremanId, companyId)
	if err != nil {
		return transactionworkshopentities.WorkOrderAllocationDetail{}, err
	}

	utils.RefreshCaches(ctx, "assign_technician")

	return entity, nil
}
