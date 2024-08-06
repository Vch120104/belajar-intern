package transactionjpcbserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BayMasterServiceImpl struct {
	BayMasterRepository transactionjpcbrepository.BayMasterRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func StartBayService(BayRepository transactionjpcbrepository.BayMasterRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.BayMasterService {
	return &BayMasterServiceImpl{
		BayMasterRepository: BayRepository,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

// GetAllBayMaster implements transactionjpcbservice.BayMasterService.
func (s *BayMasterServiceImpl) GetAllBayMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.BayMasterRepository.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

// GetAllActiveBayCarWashScreen implements transactionjpcbservice.BayMasterService.
func (s *BayMasterServiceImpl) GetAllActiveBayCarWashScreen(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.BayMasterRepository.GetAllActive(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}
