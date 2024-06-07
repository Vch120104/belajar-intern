package masteritemserviceimpl

import (
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

func (s *LandedCostMasterServiceImpl) GetAllLandedCost(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{},int,int,*exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results,totalpages,totalrows, err := s.LandedCostMasterRepo.GetAllLandedCost(tx, filterCondition, pages)
	if err != nil {
		return results,0,0,err
	}
	return results,totalpages,totalrows,nil
}

func (s *LandedCostMasterServiceImpl) GetByIdLandedCost(id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.LandedCostMasterRepo.GetByIdLandedCost(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LandedCostMasterServiceImpl) SaveLandedCost(req masteritempayloads.LandedCostMasterRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.SaveLandedCost(tx, req)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) DeactivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.DeactivateLandedCostmaster(tx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) ActivateLandedCostMaster(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.LandedCostMasterRepo.ActivateLandedCostMaster(tx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *LandedCostMasterServiceImpl) UpdateLandedCostMaster(id int,req masteritempayloads.LandedCostMasterUpdateRequest)(bool,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err :=s.LandedCostMasterRepo.UpdateLandedCostMaster(tx,id,req)
	if err != nil{
		return false,err
	}
	return result,nil
}