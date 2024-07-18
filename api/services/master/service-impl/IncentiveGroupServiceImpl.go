package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
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

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.IncentiveGroupRepository.GetAllIncentiveGroup(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroupIsActive() ([]masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.IncentiveGroupRepository.GetAllIncentiveGroupIsActive(tx)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *IncentiveGroupServiceImpl) GetIncentiveGroupById(id int) (masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *IncentiveGroupServiceImpl) SaveIncentiveGroup(req masterpayloads.IncentiveGroupResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.IncentiveGroupRepository.SaveIncentiveGroup(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupServiceImpl) ChangeStatusIncentiveGroup(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)

	if err != nil {
		// panic(exceptions.NewNotFoundError(err.Error()))
		return false, err
	}

	results, err := s.IncentiveGroupRepository.ChangeStatusIncentiveGroup(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *IncentiveGroupServiceImpl) UpdateIncentiveGroup(req masterpayloads.UpdateIncentiveGroupRequest, id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.IncentiveGroupRepository.UpdateIncentiveGroup(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroupDropDown() ([]masterpayloads.IncentiveGroupDropDown, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.IncentiveGroupRepository.GetAllIncentiveGroupDropDown(tx)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return result, err
	}
	return result, nil
}
