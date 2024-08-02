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
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WarrantyFreeServiceServiceImpl) GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.warrantyFreeServiceRepo.SaveWarrantyFreeService(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.WarrantyFreeService{}, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) ChangeStatusWarrantyFreeService(Id int) (masterpayloads.WarrantyFreeServicePatchResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)

	if err != nil {
		return masterpayloads.WarrantyFreeServicePatchResponse{}, err
	}

	results, err := s.warrantyFreeServiceRepo.ChangeStatusWarrantyFreeService(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterpayloads.WarrantyFreeServicePatchResponse{}, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) UpdateWarrantyFreeService(req masterentities.WarrantyFreeService, id int)(masterentities.WarrantyFreeService,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.warrantyFreeServiceRepo.UpdateWarrantyFreeService(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return masterentities.WarrantyFreeService{},err
	}

	return result, nil
}