package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupServiceImpl struct {
	IncentiveGroupRepository masterrepository.IncentiveGroupRepository
	DB                       *gorm.DB
}

func StartIncentiveGroupService(IncentiveGroupRepository masterrepository.IncentiveGroupRepository, db *gorm.DB) masterservice.IncentiveGroupService {
	return &IncentiveGroupServiceImpl{
		IncentiveGroupRepository: IncentiveGroupRepository,
		DB:                       db,
	}
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupRepository.GetAllIncentiveGroup(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveGroupServiceImpl) GetAllIncentiveGroupIsActive() []masterpayloads.IncentiveGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.IncentiveGroupRepository.GetAllIncentiveGroupIsActive(tx)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}

	return result
}

func (s *IncentiveGroupServiceImpl) GetIncentiveGroupById(id int) masterpayloads.IncentiveGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveGroupServiceImpl) SaveIncentiveGroup(req masterpayloads.IncentiveGroupResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.IncentiveGroupId != 0 {
		_, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, req.IncentiveGroupId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.IncentiveGroupRepository.SaveIncentiveGroup(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return results
}

func (s *IncentiveGroupServiceImpl) ChangeStatusIncentiveGroup(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.IncentiveGroupRepository.GetIncentiveGroupById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.IncentiveGroupRepository.ChangeStatusIncentiveGroup(tx, id)
	if err != nil {
		return results
	}
	return true
}
