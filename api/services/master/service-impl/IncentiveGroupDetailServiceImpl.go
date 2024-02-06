package masterserviceimpl

import (
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupDetailImpl struct {
	IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository
	DB                 *gorm.DB
}

func StartIncentiveGroupDetailService(IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository, db *gorm.DB) masterservice.IncentiveGroupDetailService {
	return &IncentiveGroupDetailImpl{
		IncentiveGroupDetailRepository: IncentiveGroupDetailRepository,
		DB:                 db,
	}
}

func (s *IncentiveGroupDetailImpl) GetIncentiveGroupDetailById(id int) masterpayloads.IncentiveGroupDetailResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveGroupDetailImpl) GetAllIncentiveGroupDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.GetAllIncentiveGroupDetail(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveGroupDetailImpl) SaveIncentiveGroupDetail(id int, req masterpayloads.IncentiveGroupDetailResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.SaveIncentiveGroupDetail(tx, id, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
