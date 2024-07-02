package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IncentiveMasterServiceImpl struct {
	IncentiveMasterRepo masterrepository.IncentiveMasterRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func StartIncentiveMasterService(IncentiveMasterRepo masterrepository.IncentiveMasterRepository, db *gorm.DB, redisClient *redis.Client) masterservice.IncentiveMasterService {
	return &IncentiveMasterServiceImpl{
		IncentiveMasterRepo: IncentiveMasterRepo,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

func (s *IncentiveMasterServiceImpl) GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.IncentiveMasterRepo.GetAllIncentiveMaster(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *IncentiveMasterServiceImpl) GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.IncentiveMasterRepo.GetIncentiveMasterById(tx, id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *IncentiveMasterServiceImpl) SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.IncentiveMasterRepo.SaveIncentiveMaster(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *IncentiveMasterServiceImpl) ChangeStatusIncentiveMaster(Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.IncentiveMasterRepo.ChangeStatusIncentiveMaster(tx, Id)
	if err != nil {
		return masterentities.IncentiveMaster{}, err
	}

	defer helper.CommitOrRollback(tx, err)
	return entity, nil
}
