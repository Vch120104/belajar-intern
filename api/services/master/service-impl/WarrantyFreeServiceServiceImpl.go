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

type WarrantyFreeServiceServiceImpl struct {
	warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository
	DB                      *gorm.DB
	RedisClient             *redis.Client // Redis client
}

func StartWarrantyFreeServiceService(warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository, db *gorm.DB, redisClient *redis.Client) masterservice.WarrantyFreeServiceService {
	return &WarrantyFreeServiceServiceImpl{
		warrantyFreeServiceRepo: warrantyFreeServiceRepo,
		DB:                      db,
		RedisClient:             redisClient,
	}
}

func (s *WarrantyFreeServiceServiceImpl) GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.warrantyFreeServiceRepo.GetAllWarrantyFreeService(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *WarrantyFreeServiceServiceImpl) GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)

	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.WarrantyFreeServicesId != 0 {
		_, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, req.WarrantyFreeServicesId)
		if err != nil {
			return false, err
		}
	}

	results, err := s.warrantyFreeServiceRepo.SaveWarrantyFreeService(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) ChangeStatusWarrantyFreeService(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.warrantyFreeServiceRepo.ChangeStatusWarrantyFreeService(tx, Id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}
