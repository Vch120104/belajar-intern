package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationServiceImpl struct {
	warehouseLocationRepo masterwarehouserepository.WarehouseLocationRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func OpenWarehouseLocationService(warehouseLocation masterwarehouserepository.WarehouseLocationRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseLocationService {
	return &WarehouseLocationServiceImpl{
		warehouseLocationRepo: warehouseLocation,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

func (s *WarehouseLocationServiceImpl) Save(request masterwarehousepayloads.GetWarehouseLocationResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.WarehouseLocationId != 0 {
		_, err := s.warehouseLocationRepo.GetById(tx, request.WarehouseLocationId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.warehouseLocationRepo.Save(tx, request)

	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return save, err
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return get, err
	}
	defer helper.CommitOrRollback(tx, err)
	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseLocationRepo.GetAll(tx, filter, pages)

	if err != nil {
		return get, err
	}
	defer helper.CommitOrRollback(tx, err)
	return get, nil
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)

	if err != nil {
		return change_status, err
	}
	defer helper.CommitOrRollback(tx, err)
	return change_status, nil
}
