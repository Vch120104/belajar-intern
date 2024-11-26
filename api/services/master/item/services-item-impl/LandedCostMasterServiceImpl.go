package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
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

func (s *LandedCostMasterServiceImpl) GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.LandedCostMasterRepo.GetAllLandedCost(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LandedCostMasterServiceImpl) GetByIdLandedCost(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.LandedCostMasterRepo.GetByIdLandedCost(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LandedCostMasterServiceImpl) SaveLandedCost(req masteritempayloads.LandedCostMasterRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.LandedCostMasterRepo.SaveLandedCost(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.LandedCost{}, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) DeactivateLandedCostMaster(id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.LandedCostMasterRepo.DeactivateLandedCostmaster(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) ActivateLandedCostMaster(id string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.LandedCostMasterRepo.ActivateLandedCostMaster(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) UpdateLandedCostMaster(id int, req masteritempayloads.LandedCostMasterUpdateRequest) (masteritementities.LandedCost, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.LandedCostMasterRepo.UpdateLandedCostMaster(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.LandedCost{}, err
	}
	return result, nil
}
