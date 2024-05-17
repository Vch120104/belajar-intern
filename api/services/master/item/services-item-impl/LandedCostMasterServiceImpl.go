package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
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

func (s *LandedCostMasterServiceImpl) GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.LandedCostMasterRepo.GetAllLandedCost(tx, filterCondition, pages)
	if err != nil {
		return results,err
	}
	return results,nil
}

func (s *LandedCostMasterServiceImpl) GetByIdLandedCost(id int) (masteritempayloads.LandedCostMasterPayloads,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.LandedCostMasterRepo.GetByIdLandedCost(tx, id)
	if err != nil {
		return results,err
	}
	return results,nil
}

func (s *LandedCostMasterServiceImpl) SaveLandedCost(req masteritempayloads.LandedCostMasterPayloads) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.SaveLandedCost(tx, req)
	if err != nil {
		return false,err
	}
	return result,nil
}

func (s *LandedCostMasterServiceImpl) DeactivateLandedCostMaster(id string) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.DeactivateLandedCostmaster(tx, id)
	if err != nil {
		return false,err
	}
	return result,nil
}

func (s *LandedCostMasterServiceImpl) ActivateLandedCostMaster(id string) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.ActivateLandedCostMaster(tx, id)
	if err != nil {
		return false,err
	}
	return result,nil
}
