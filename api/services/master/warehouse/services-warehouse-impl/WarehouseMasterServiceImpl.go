package masterwarehouseserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

	// "log"

	// "after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WarehouseMasterServiceImpl struct {
	warehouseMasterRepo masterwarehouserepository.WarehouseMasterRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func OpenWarehouseMasterService(warehouseMaster masterwarehouserepository.WarehouseMasterRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseMasterService {
	return &WarehouseMasterServiceImpl{
		warehouseMasterRepo: warehouseMaster,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

func (s *WarehouseMasterServiceImpl) Save(request masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	save, err := s.warehouseMasterRepo.Save(tx, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WarehouseMasterServiceImpl) GetById(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetById(tx, warehouseId)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.DropdownWarehouse(tx)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetAllIsActive(tx)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetWarehouseWithMultiId(tx, MultiIds)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetAll(tx, request, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetWarehouseMasterByCode(tx, Code)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) ChangeStatus(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	change_status, err := s.warehouseMasterRepo.ChangeStatus(tx, warehouseId)

	if err != nil {
		return change_status, err
	}

	return change_status, nil
}
