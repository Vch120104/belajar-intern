package transactionworkshopserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkOrderServiceImpl struct {
	structWorkOrderRepo transactionworkshoprepository.WorkOrderRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func OpenWorkOrderServiceImpl(WorkOrderRepo transactionworkshoprepository.WorkOrderRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.WorkOrderService {
	return &WorkOrderServiceImpl{
		structWorkOrderRepo: WorkOrderRepo,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

func (s *WorkOrderServiceImpl) WithTrx(Trxhandle *gorm.DB) transactionworkshopservice.WorkOrderService {
	s.structWorkOrderRepo = s.structWorkOrderRepo.WithTrx(Trxhandle)
	return s
}

func (s *WorkOrderServiceImpl) Save(request transactionworkshoppayloads.WorkOrderRequest) (bool, error) {
	save, err := s.structWorkOrderRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WorkOrderServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}
