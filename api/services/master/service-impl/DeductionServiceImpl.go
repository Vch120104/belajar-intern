package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// "after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DeductionServiceImpl struct {
	deductionrepo masterrepository.DeductionRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartDeductionService(deductionRepo masterrepository.DeductionRepository, db *gorm.DB, redisClient *redis.Client) masterservice.DeductionService {
	return &DeductionServiceImpl{
		deductionrepo: deductionRepo,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *DeductionServiceImpl) GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	// Generate key for caching
	cacheKey := "deduction:all"

	// Check if data is available in cache
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// If data found in cache, return it
		var result pagination.Pagination
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	}

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetAllDeduction(tx, filterCondition, pages)
	if dbErr != nil {
		// Handle error from the database operation
		return pagination.Pagination{}, dbErr
	}

	// Store data in cache for future use
	jsonData, _ := json.Marshal(result)
	if err := s.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
		// Log the error or handle it appropriately
		log.Println("Error storing data in cache:", err)
		// Atau lakukan penanganan kesalahan yang sesuai
	}

	defer helper.CommitOrRollback(tx, dbErr)
	return result, nil
}

func (s *DeductionServiceImpl) GetByIdDeductionDetail(Id int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse) {
	ctx := context.Background() // Initialize context

	// Generate key for caching
	cacheKey := fmt.Sprintf("deduction:detail:%d", Id)

	// Check if data is available in cache
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// If data found in cache, return it
		var result masterpayloads.DeductionDetailResponse
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	}

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetByIdDeductionDetail(tx, Id)
	if dbErr != nil {
		// Handle error
		return result, dbErr // Return the existing BaseErrorResponse
	}

	// Store data in cache for future use
	jsonData, _ := json.Marshal(result)
	if err := s.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
		// Log or handle error
		log.Println("Error storing data in cache:", err)
	}
	defer helper.CommitOrRollback(tx, dbErr)
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionList(req masterpayloads.DeductionListResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.deductionrepo.SaveDeductionList(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionDetail(req masterpayloads.DeductionDetailResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.deductionrepo.SaveDeductionDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetDeductionById(Id int) (masterpayloads.DeductionListResponse, *exceptions.BaseErrorResponse) {
	ctx := context.Background() // Inisialisasi context

	// Generate key for caching
	cacheKey := fmt.Sprintf("deduction:id:%d", Id)

	// Check if data is available in cache
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// If data found in cache, return it
		var result masterpayloads.DeductionListResponse
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	}

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetDeductionById(tx, Id)
	defer helper.CommitOrRollback(tx, dbErr)
	if dbErr != nil {
		// Handle error
		return masterpayloads.DeductionListResponse{}, dbErr
	}

	// Store data in cache for future use
	jsonData, _ := json.Marshal(result)
	if err := s.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err(); err != nil {
		// Log or handle error
		log.Println("Error storing data in cache:", err)
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetAllDeductionDetail(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	detail_result, detail_err := s.deductionrepo.GetAllDeductionDetail(tx, pages, Id)
	defer helper.CommitOrRollback(tx, detail_err)

	if detail_err != nil {
		return detail_result, detail_err
	}
	return detail_result, nil
}

func (s *DeductionServiceImpl) ChangeStatusDeduction(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.deductionrepo.GetDeductionById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.deductionrepo.ChangeStatusDeduction(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return true, nil
}
