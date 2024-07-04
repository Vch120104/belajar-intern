package masterserviceimpl

import (
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

type ForecastMasterServiceImpl struct {
	ForecastMasterRepo masterrepository.ForecastMasterRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func StartForecastMasterService(ForecastMasterRepo masterrepository.ForecastMasterRepository, db *gorm.DB, redisClient *redis.Client) masterservice.ForecastMasterService {
	return &ForecastMasterServiceImpl{
		ForecastMasterRepo: ForecastMasterRepo,
		DB:                 db,
		RedisClient:        redisClient,
	}
}

func (s *ForecastMasterServiceImpl) GetForecastMasterById(id int) (masterpayloads.ForecastMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ForecastMasterRepo.GetForecastMasterById(tx, id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ForecastMasterServiceImpl) SaveForecastMaster(req masterpayloads.ForecastMasterResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.ForecastMasterId != 0 {
		_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, req.ForecastMasterId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.ForecastMasterRepo.SaveForecastMaster(tx, req)

	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ForecastMasterServiceImpl) ChangeStatusForecastMaster(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.ForecastMasterRepo.ChangeStatusForecastMaster(tx, Id)
	if err != nil {
		return results, nil
	}
	defer helper.CommitOrRollback(tx, err)
	return true, nil
}

func (s *ForecastMasterServiceImpl) GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ForecastMasterRepo.GetAllForecastMaster(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *ForecastMasterServiceImpl) UpdateForecastMaster(req masterpayloads.ForecastMasterResponse, id int)(bool,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,err := s.ForecastMasterRepo.UpdateForecastMaster(tx,req,id)
	if err != nil{
		return false,err
	}
	return result,nil
}