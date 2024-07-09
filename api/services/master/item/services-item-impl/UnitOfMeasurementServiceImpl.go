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

type UnitOfMeasurementServiceImpl struct {
	unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func StartUnitOfMeasurementService(unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.UnitOfMeasurementService {
	return &UnitOfMeasurementServiceImpl{
		unitOfMeasurementRepo: unitOfMeasurementRepo,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurementIsActive(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementById(id int) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementByCode(tx, Code)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurement(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) ChangeStatusUnitOfMeasurement(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.unitOfMeasurementRepo.ChangeStatusUnitOfMeasurement(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.UomId != 0 {
		_, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, req.UomId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.unitOfMeasurementRepo.SaveUnitOfMeasurement(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
