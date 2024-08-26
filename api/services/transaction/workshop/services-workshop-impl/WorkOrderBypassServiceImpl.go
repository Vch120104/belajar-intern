package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkOrderBypassServiceImpl struct {
	structWorkOrderBypassRepo transactionworkshoprepository.WorkOrderBypassRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func OpenWorkOrderBypassServiceImpl(WorkOrderBypassRepo transactionworkshoprepository.WorkOrderBypassRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.WorkOrderBypassService {
	return &WorkOrderBypassServiceImpl{
		structWorkOrderBypassRepo: WorkOrderBypassRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *WorkOrderBypassServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.structWorkOrderBypassRepo.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *WorkOrderBypassServiceImpl) GetById(id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	workOrder, err := s.structWorkOrderBypassRepo.GetById(tx, id)
	if err != nil {
		return workOrder, err
	}

	return workOrder, nil
}

func (s *WorkOrderBypassServiceImpl) Bypass(request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshopentities.WorkOrderQualityControl, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	workOrder, err := s.structWorkOrderBypassRepo.Bypass(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return workOrder, err
	}

	// if workOrder.Status != "BYPASS" {
	// 	return workOrder, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusUnprocessableEntity,
	// 		Err:        fmt.Errorf("Work order status is not BYPASS"),
	// 	}
	// }

	// workOrder.Status = "QUALITY_CONTROL"
	// if err := s.structWorkOrderBypassRepo.UpdateWorkOrder(tx, workOrder); err != nil {
	// 	return workOrder, err
	// }

	// // Create quality control
	// qualityControl := transactionworkshopentities.WorkOrderQualityControl{
	// 	WorkOrderID: workOrder.ID,
	// 	Status:      "OPEN",
	// }

	// if err := s.structWorkOrderBypassRepo.CreateQualityControl(tx, qualityControl); err != nil {
	// 	return workOrder, err
	// }

	return workOrder, nil
}
