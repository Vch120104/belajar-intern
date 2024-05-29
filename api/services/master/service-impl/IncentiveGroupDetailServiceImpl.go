package masterserviceimpl

import (
	// "after-sales/api/exceptions"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IncentiveGroupDetailImpl struct {
	IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository
	DB                             *gorm.DB
	RedisClient                    *redis.Client // Redis client
}

func StartIncentiveGroupDetailService(IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository, db *gorm.DB, redisClient *redis.Client) masterservice.IncentiveGroupDetailService {
	return &IncentiveGroupDetailImpl{
		IncentiveGroupDetailRepository: IncentiveGroupDetailRepository,
		DB:                             db,
		RedisClient:                    redisClient,
	}
}

func (s *IncentiveGroupDetailImpl) GetAllIncentiveGroupDetail(headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.GetAllIncentiveGroupDetail(tx, headerId, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupDetailImpl) GetIncentiveGroupDetailById(id int) (masterpayloads.IncentiveGroupDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupDetailImpl) SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.IncentiveGroupDetailId != 0 {
		_, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, req.IncentiveGroupDetailId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.IncentiveGroupDetailRepository.SaveIncentiveGroupDetail(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}
