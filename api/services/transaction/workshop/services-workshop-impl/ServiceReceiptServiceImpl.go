package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
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

type ServiceReceiptServiceImpl struct {
	ServiceReceiptRepository transactionworkshoprepository.ServiceReceiptRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenServiceReceiptServiceImpl(ServiceReceiptRepo transactionworkshoprepository.ServiceReceiptRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceReceiptService {
	return &ServiceReceiptServiceImpl{
		ServiceReceiptRepository: ServiceReceiptRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *ServiceReceiptServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.ServiceReceiptRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil
}

func (s *ServiceReceiptServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse) {

	cacheKey := utils.GenerateCacheKeyIds("service_receipt_id", id)

	ctx := context.Background()
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var result transactionworkshoppayloads.ServiceReceiptResponse
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		return result, nil
	} else if err != redis.Nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.ServiceReceiptRepository.GetById(tx, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		s.RedisClient.Set(ctx, cacheKey, cacheData, utils.CacheExpiration)
	}

	return result, nil
}

func (s *ServiceReceiptServiceImpl) Save(id int, request transactionworkshoppayloads.ServiceReceiptSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	ctx := context.Background()

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	save, err := s.ServiceReceiptRepository.Save(tx, id, request)
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, err
	}

	utils.RefreshCaches(ctx, "service_receipt")

	return save, nil
}
