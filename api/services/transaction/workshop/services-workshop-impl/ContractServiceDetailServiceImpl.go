package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ContractServiceDetailServiceImpl struct {
	ContractServiceDetailRepository transactionworkshoprepository.ContractServiceDetailRepository
	DB                              *gorm.DB
	RedisClient                     *redis.Client
}

// Fungsi pembuka untuk membuat instance ContractServiceServiceImpl
func OpenContractServiceDetailServiceImpl(ContractServiceDetailRepo transactionworkshoprepository.ContractServiceDetailRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ContractServiceDetailService {
	return &ContractServiceDetailServiceImpl{
		ContractServiceDetailRepository: ContractServiceDetailRepo,
		DB:                              db,
		RedisClient:                     redisClient,
	}
}

// GetAllDetail implements transactionworkshopservice.ContractServiceDetailService.
// GetAllDetail implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) GetAllDetail(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// Memulai transaksi
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Memanggil repository untuk mengambil data ContractServiceDetail berdasarkan contract_service_system_number
	results, totalPages, totalRows, repoErr := s.ContractServiceDetailRepository.GetAllDetail(tx, Id, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	// Menggunakan pagination untuk hasil
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// GetById implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) GetById(Id int) (transactionworkshoppayloads.ContractServiceIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ContractServiceDetailRepository.GetById(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
