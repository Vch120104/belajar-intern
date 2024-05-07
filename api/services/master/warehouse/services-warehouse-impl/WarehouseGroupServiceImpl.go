package masterwarehouseserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type WarehouseGroupServiceImpl struct {
	warehouseGroupRepo masterwarehouserepository.WarehouseGroupRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func OpenWarehouseGroupService(warehouseGroup masterwarehouserepository.WarehouseGroupRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseGroupService {
	return &WarehouseGroupServiceImpl{
		warehouseGroupRepo: warehouseGroup,
		DB:                 db,
		RedisClient:        redisClient,
	}
}

func (s *WarehouseGroupServiceImpl) SaveWarehouseGroup(request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.WarehouseGroupId != 0 {
		_, err := s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, request.WarehouseGroupId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.warehouseGroupRepo.SaveWarehouseGroup(tx, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WarehouseGroupServiceImpl) GetByIdWarehouseGroup(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseGroupServiceImpl) GetAllWarehouseGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseGroupRepo.GetAllWarehouseGroup(tx, filterCondition, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *WarehouseGroupServiceImpl) ChangeStatusWarehouseGroup(warehouseGroupId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseGroupRepo.ChangeStatusWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return change_status, err
	}

	return change_status, nil
}
