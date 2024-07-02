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

type OperationCodeServiceImpl struct {
	operationCodeRepo masteroperationrepository.OperationCodeRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartOperationCodeService(operationCodeRepo masteroperationrepository.OperationCodeRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationCodeService {
	return &OperationCodeServiceImpl{
		operationCodeRepo: operationCodeRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *OperationCodeServiceImpl) GetAllOperationCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationCodeRepo.GetAllOperationCode(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationCodeServiceImpl) GetOperationCodeById(id int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationCodeRepo.GetOperationCodeById(tx, id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationCodeServiceImpl) GetOperationCodeByCode(code string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationCodeRepo.GetOperationCodeByCode(tx, code)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *OperationCodeServiceImpl) SaveOperationCode(req masteroperationpayloads.OperationCodeSave) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.operationCodeRepo.SaveOperationCode(tx, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *OperationCodeServiceImpl) ChangeStatusOperationCode(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Statement.DB.Begin()
	result, err := s.operationCodeRepo.ChangeStatusItemSubstitute(tx, id)

	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}
