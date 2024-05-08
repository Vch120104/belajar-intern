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

type DiscountServiceImpl struct {
	discountRepo masterrepository.DiscountRepository
	DB           *gorm.DB
	RedisClient  *redis.Client // Redis client
}

func StartDiscountService(discountRepo masterrepository.DiscountRepository, db *gorm.DB, redisClient *redis.Client) masterservice.DiscountService {
	return &DiscountServiceImpl{
		discountRepo: discountRepo,
		DB:           db,
		RedisClient:  redisClient,
	}
}

func (s *DiscountServiceImpl) GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetAllDiscountIsActive(tx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountById(id int) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetDiscountById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetDiscountByCode(tx, Code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetAllDiscount(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) ChangeStatusDiscount(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.discountRepo.GetDiscountById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.discountRepo.ChangeStatusDiscount(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *DiscountServiceImpl) SaveDiscount(req masterpayloads.DiscountResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.DiscountCodeId != 0 {
		_, err := s.discountRepo.GetDiscountById(tx, req.DiscountCodeId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.discountRepo.SaveDiscount(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}
