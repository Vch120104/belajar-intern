package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DiscountPercentServiceImpl struct {
	discountPercentRepo masteritemrepository.DiscountPercentRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func StartDiscountPercentService(discountPercentRepo masteritemrepository.DiscountPercentRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.DiscountPercentService {
	return &DiscountPercentServiceImpl{
		discountPercentRepo: discountPercentRepo,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

func (s *DiscountPercentServiceImpl) GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.discountPercentRepo.GetAllDiscountPercent(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *DiscountPercentServiceImpl) GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *DiscountPercentServiceImpl) SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	if req.DiscountPercentId != 0 {
		_, err := s.discountPercentRepo.GetDiscountPercentById(tx, req.DiscountPercentId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.discountPercentRepo.SaveDiscountPercent(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *DiscountPercentServiceImpl) ChangeStatusDiscountPercent(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.discountPercentRepo.ChangeStatusDiscountPercent(tx, Id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return true, nil
}
