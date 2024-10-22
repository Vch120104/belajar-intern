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

// Fungsi pembuka untuk membuat instance ContractServiceServiceImpl
func OpenContractServiceServiceImpl(contractServiceRepo transactionworkshoprepository.ContractServiceRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ContractServiceService {
	return &ContractServiceServiceImpl{
		ContractServiceRepository: contractServiceRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

// GetAll implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// Memulai transaksi
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Mengambil data dari repository
	results, totalPages, totalRows, repoErr := s.ContractServiceRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	// Menggunakan pagination untuk hasil
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// GetById implements transactionworkshopservice.ContractServiceService.
func (s *ContractServiceServiceImpl) GetById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse) {
	// Membuat cache key untuk menyimpan data di Redis
	cacheKey := utils.GenerateCacheKeyIds("contract_service_system_number", Id)
	ctx := context.Background()

	// Cek apakah data sudah ada di cache Redis
	cachedData, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Jika data ditemukan di cache, unmarshal dan kembalikan sebagai response
		var result transactionworkshoppayloads.ContractServiceResponseId
		if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
			return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to unmarshal cached data",
				Err:        err,
			}
		}
		// Data dari Redis ditemukan, kembalikan hasil
		return result, nil
	} else if err != redis.Nil {
		// Jika ada error selain key tidak ditemukan, kembalikan error
		fmt.Println("Redis error:", err)
		return transactionworkshoppayloads.ContractServiceResponseId{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Redis error",
			Err:        err,
		}
	}

	// Mulai transaksi database
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika ada panic
		}
	}()

	// Mengambil data dari repository
	result, repoErr := s.ContractServiceRepository.GetById(tx, Id, filterCondition, pages)
	if repoErr != nil {
		tx.Rollback() // Rollback jika ada error pada repository
		return result, repoErr
	}

	// Commit transaksi jika berhasil
	tx.Commit()

	// Cache data hasil dari repository ke Redis
	cacheData, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		fmt.Println("Failed to marshal result for caching:", marshalErr)
	} else {
		if err := s.RedisClient.Set(ctx, cacheKey, cacheData, utils.CacheExpiration).Err(); err != nil {
			fmt.Println("Failed to set cache:", err)
		}
	}

	// Kembalikan hasil
	return result, nil
}
