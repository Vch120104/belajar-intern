package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"errors"
	"net/http"

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

// UpdateDiscount implements masterservice.DiscountService.
func (s *DiscountServiceImpl) UpdateDiscount(id int, req masterpayloads.DiscountUpdate) (bool, *exceptions.BaseErrorResponse) {

	//check id
	res, _ := s.GetDiscountById(id)

	if res == (masterpayloads.DiscountResponse{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("id not found"),
			Message:    "Id not found",
		}
	}

	tx := s.DB.Begin()
	results, err := s.discountRepo.UpdateDiscount(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.discountRepo.GetAllDiscountIsActive(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountById(id int) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.discountRepo.GetDiscountById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.discountRepo.GetDiscountByCode(tx, Code)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.discountRepo.GetAllDiscount(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) ChangeStatusDiscount(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.discountRepo.GetDiscountById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.discountRepo.ChangeStatusDiscount(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *DiscountServiceImpl) SaveDiscount(req masterpayloads.DiscountResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.DiscountCodeId != 0 {
		_, err := s.discountRepo.GetDiscountById(tx, req.DiscountCodeId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.discountRepo.SaveDiscount(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}
