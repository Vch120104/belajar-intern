package masteroperationserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OperationEntriesServiceImpl struct {
	operationEntriesRepo masteroperationrepository.OperationEntriesRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func StartOperationEntriesService(operationEntriesRepo masteroperationrepository.OperationEntriesRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationEntriesService {
	return &OperationEntriesServiceImpl{
		operationEntriesRepo: operationEntriesRepo,
		DB:                   db,
		RedisClient:          redisClient,
	}
}

func (s *OperationEntriesServiceImpl) GetAllOperationEntries(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationEntriesRepo.GetAllOperationEntries(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesById(id int) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationEntriesRepo.GetOperationEntriesById(tx, id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesName(request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationEntriesRepo.GetOperationEntriesName(tx, request)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationEntriesServiceImpl) SaveOperationEntries(req masteroperationpayloads.OperationEntriesResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.OperationEntriesId != 0 {
		_, err := s.operationEntriesRepo.GetOperationEntriesById(tx, req.OperationEntriesId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationEntriesRepo.SaveOperationEntries(tx, req)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationEntriesServiceImpl) ChangeStatusOperationEntries(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationEntriesRepo.ChangeStatusOperationEntries(tx, Id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return true, nil
}
