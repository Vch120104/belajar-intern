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

type OperationKeyServiceImpl struct {
	operationKeyRepo masteroperationrepository.OperationKeyRepository
	DB               *gorm.DB
	RedisClient      *redis.Client // Redis client
}

func StartOperationKeyService(operationKeyRepo masteroperationrepository.OperationKeyRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationKeyService {
	return &OperationKeyServiceImpl{
		operationKeyRepo: operationKeyRepo,
		DB:               db,
		RedisClient:      redisClient,
	}
}

func (s *OperationKeyServiceImpl) GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationKeyRepo.GetAllOperationKeyList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyById(id int) (masteroperationpayloads.OperationkeyListResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationKeyRepo.GetOperationKeyById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyName(req masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationKeyRepo.GetOperationKeyName(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err

	}
	return results, nil
}

func (s *OperationKeyServiceImpl) SaveOperationKey(req masteroperationpayloads.OperationKeyResponse) (bool, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()

	if req.OperationKeyId != 0 {
		_, err := s.operationKeyRepo.GetOperationKeyById(tx, req.OperationKeyId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationKeyRepo.SaveOperationKey(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) ChangeStatusOperationKey(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationKeyRepo.ChangeStatusOperationKey(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
