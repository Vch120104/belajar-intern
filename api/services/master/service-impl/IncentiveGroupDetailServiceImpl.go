package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type IncentiveGroupDetailImpl struct {
	IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository
	DB                             *gorm.DB
}

func StartIncentiveGroupDetailService(IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository, db *gorm.DB) masterservice.IncentiveGroupDetailService {
	return &IncentiveGroupDetailImpl{
		IncentiveGroupDetailRepository: IncentiveGroupDetailRepository,
		DB:                             db,
	}
}

func (s *IncentiveGroupDetailImpl) GetAllIncentiveGroupDetail(headerId int, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.GetAllIncentiveGroupDetail(tx, headerId, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
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

func (s *IncentiveGroupDetailImpl) SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.IncentiveGroupDetailId != 0 {
		_, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, req.IncentiveGroupDetailId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.IncentiveGroupDetailRepository.SaveIncentiveGroupDetail(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
