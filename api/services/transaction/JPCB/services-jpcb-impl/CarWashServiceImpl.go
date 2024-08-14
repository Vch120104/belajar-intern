package transactionjpcbserviceimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CarWashServiceImpl struct {
	CarWashRepository transactionjpcbrepository.CarWashRepository
	DB                *gorm.DB
	RedisClient       *redis.Client
}

func NewCarWashServiceImpl(CarWashRepository transactionjpcbrepository.CarWashRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.CarWashService {
	return &CarWashServiceImpl{
		CarWashRepository: CarWashRepository,
		DB:                db,
		RedisClient:       redisClient,
	}
}

// GetAll implements transactionjpcbservice.CarWashService.
func (s *CarWashServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.CarWashRepository.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

// UpdatePriority implements transactionjpcbservice.CarWashService.
func (s *CarWashServiceImpl) UpdatePriority(workOrderSystemNumber int, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.UpdatePriority(tx, workOrderSystemNumber, carWashPriorityId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionjpcbentities.CarWash{}, err
	}
	return result, nil
}
