package masterwarehouseserviceimpl

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WarehouseCostingTypeServiceImpl struct {
	costingTypeRepo masterwarehouserepository.WarehouseCostingTypeRepository
	DB              *gorm.DB
	RedisClient     *redis.Client // Redis client
}

func NewWarehouseCostingTypeServiceImpl(costingTypeRepo masterwarehouserepository.WarehouseCostingTypeRepository, DB *gorm.DB, RedisClient *redis.Client) masterwarehouseservice.WarehouseCostingTypeService {
	return &WarehouseCostingTypeServiceImpl{costingTypeRepo: costingTypeRepo, DB: DB, RedisClient: RedisClient}
}
func (service *WarehouseCostingTypeServiceImpl) GetByCodeWarehouseCostingType(CostingCode string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	get, err := service.costingTypeRepo.GetByCodeWarehouseCostingType(tx, CostingCode)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return get, err
	}
	return get, nil
}
