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

func (s *WorkOrderServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.structWorkOrderRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WorkOrderServiceImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.New(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *WorkOrderServiceImpl) GetById(id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.structWorkOrderRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WorkOrderServiceImpl) Save(request transactionworkshoppayloads.WorkOrderRequest) (bool, error) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	save, err := s.structWorkOrderRepo.Save(request)
	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WorkOrderServiceImpl) Submit(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.Submit(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) Void(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.Void(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkOrderServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptionsss_test.BaseErrorResponse {
	defer helper.CommitOrRollback(tx)
	err := s.structWorkOrderRepo.CloseOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}
