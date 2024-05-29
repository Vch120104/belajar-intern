package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

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
	defer helper.CommitOrRollback(tx)

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

	return save, err
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetAll(tx, request, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)

	if err != nil {
		return change_status, err
	}

	return change_status, nil
}
