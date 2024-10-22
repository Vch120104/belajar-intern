package transactionworkshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

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
