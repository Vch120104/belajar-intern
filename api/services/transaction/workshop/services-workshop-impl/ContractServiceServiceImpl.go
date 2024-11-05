package transactionworkshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ContractServiceServiceImpl struct {
	ContractServiceRepository transactionworkshoprepository.ContractServiceRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client
}

func OpenContractServiceServiceImpl(contractServiceRepo transactionworkshoprepository.ContractServiceRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ContractServiceService {
	return &ContractServiceServiceImpl{
		ContractServiceRepository: contractServiceRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

// GetAll implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.ContractServiceRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// GetById implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) GetById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse) {
	cacheKey := utils.GenerateCacheKeyIds("contract_service_system_number", Id)
	ctx := context.Background()

	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result transactionworkshoppayloads.ContractServiceResponseId
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to unmarshal cached data",
				Err:        err,
			}
		}
		return result, nil
	} else if err != redis.Nil {
		fmt.Println("Redis error:", err)
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Redis error",
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result, repoErr := s.ContractServiceRepository.GetById(tx, Id, filterCondition, pages)
	if repoErr != nil {
		tx.Rollback()
		return result, repoErr
	}
	tx.Commit()
	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		if err := s.RedisClient.Set(ctx, cacheKey, cacheData, utils.CacheExpiration).Err(); err != nil {
			fmt.Println("Failed to set cache:", err)
		}
	}
	return result, nil
}

// Save implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) Save(payload transactionworkshoppayloads.ContractServiceInsert) (transactionworkshoppayloads.ContractServiceInsert, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, err := s.ContractServiceRepository.Save(tx, payload)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Void implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) Void(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	delete, err := s.ContractServiceRepository.Void(tx, Id)
	if err != nil {
		return false, err
	}
	return delete, nil
}

// Submit implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) Submit(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	submit, err := s.ContractServiceRepository.Submit(tx, Id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}

	return submit, nil
}
