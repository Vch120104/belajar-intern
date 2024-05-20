package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IncentiveGroupServiceImpl struct {
	IncentiveGroupRepository masterrepository.IncentiveGroupRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func StartIncentiveGroupService(IncentiveGroupRepository masterrepository.IncentiveGroupRepository, db *gorm.DB, redisClient *redis.Client) masterservice.IncentiveGroupService {
	return &IncentiveGroupServiceImpl{
		IncentiveGroupRepository: IncentiveGroupRepository,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.IncentiveGroupRepository.GetAllIncentiveGroup(tx, filterCondition, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.IncentiveGroupRepository.GetAllIncentiveGroupIsActive(tx)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *IncentiveGroupServiceImpl) GetIncentiveGroupById(id int) (masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *IncentiveGroupServiceImpl) SaveIncentiveGroup(req masterpayloads.IncentiveGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupRepository.SaveIncentiveGroup(tx, req)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *IncentiveGroupServiceImpl) ChangeStatusIncentiveGroup(id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)

	if err != nil {
		// panic(exceptions.NewNotFoundError(err.Error()))
		return false, err
	}

	results, err := s.IncentiveGroupRepository.ChangeStatusIncentiveGroup(tx, id)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *IncentiveGroupServiceImpl) UpdateIncentiveGroup(req masterpayloads.UpdateIncentiveGroupRequest, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupRepository.UpdateIncentiveGroup(tx, id, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroupDropDown() ([]masterpayloads.IncentiveGroupDropDown, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.IncentiveGroupRepository.GetAllIncentiveGroupDropDown(tx)

	if err != nil {
		return result, err
	}

	return result, nil
}
