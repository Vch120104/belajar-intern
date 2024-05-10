package masteroperationserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OperationGroupServiceImpl struct {
	operationGroupRepo masteroperationrepository.OperationGroupRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func StartOperationGroupService(operationGroupRepo masteroperationrepository.OperationGroupRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationGroupService {
	return &OperationGroupServiceImpl{
		operationGroupRepo: operationGroupRepo,
		DB:                 db,
		RedisClient:        redisClient,
	}
}

func (s *OperationGroupServiceImpl) GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.operationGroupRepo.GetAllOperationGroupIsActive(tx)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupById(id int) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationGroupRepo.GetOperationGroupById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupByCode(Code string) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationGroupRepo.GetOperationGroupByCode(tx, Code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (service *OperationGroupServiceImpl) GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	// tx := s.DB.Begin()
	// defer helper.CommitOrRollback(tx)
	// results, err := s.operationGroupRepo.GetAllOperationGroup(tx, filterCondition, pages)
	// if err != nil {
	// 	panic(exceptions.NewNotFoundError(err.Error()))
	// }
	// return results
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := service.operationGroupRepo.GetAllOperationGroup(tx, filterCondition, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *OperationGroupServiceImpl) ChangeStatusOperationGroup(oprId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.operationGroupRepo.GetOperationGroupById(tx, oprId)

	if err != nil {
		return false, err
	}

	results, err := s.operationGroupRepo.ChangeStatusOperationGroup(tx, oprId)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *OperationGroupServiceImpl) SaveOperationGroup(req masteroperationpayloads.OperationGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationGroupId != 0 {
		_, err := s.operationGroupRepo.GetOperationGroupById(tx, req.OperationGroupId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationGroupRepo.SaveOperationGroup(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}
