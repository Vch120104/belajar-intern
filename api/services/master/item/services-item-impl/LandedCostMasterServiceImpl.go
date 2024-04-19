package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LandedCostMasterServiceImpl struct {
	LandedCostMasterRepo masteritemrepository.LandedCostMasterRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func StartLandedCostMasterService(LandedCostMasterRepo masteritemrepository.LandedCostMasterRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.LandedCostMasterService {
	return &LandedCostMasterServiceImpl{
		LandedCostMasterRepo: LandedCostMasterRepo,
		DB:                   db,
		RedisClient:          redisClient,
	}
}

func (s *LandedCostMasterServiceImpl) GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.LandedCostMasterRepo.GetAllLandedCost(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return results
}

func (s *LandedCostMasterServiceImpl) GetByIdLandedCost(id int) masteritempayloads.LandedCostMasterPayloads {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.LandedCostMasterRepo.GetByIdLandedCost(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return results
}

func (s *LandedCostMasterServiceImpl) SaveLandedCost(req masteritempayloads.LandedCostMasterPayloads) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.SaveLandedCost(tx, req)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *LandedCostMasterServiceImpl) DeactivateLandedCostMaster(id string) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.DeactivateLandedCostmaster(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *LandedCostMasterServiceImpl) ActivateLandedCostMaster(id string) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.ActivateLandedCostMaster(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}
