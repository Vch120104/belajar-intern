package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VehicleHistoryServiceImpl struct {
	VehicleHistoryRepo transactionworkshoprepository.VehicleHistoryRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func NewVehicleHistoryServiceImpl(VehicleHistoryRepo transactionworkshoprepository.VehicleHistoryRepository, db *gorm.DB, redis *redis.Client) transactionworkshopservice.VehicleHistoryService {
	return &VehicleHistoryServiceImpl{
		VehicleHistoryRepo: VehicleHistoryRepo,
		DB:                 db,
		RedisClient:        redis,
	}
}
func (v *VehicleHistoryServiceImpl) GetAllVehicleHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := v.DB.Begin()
	result, err := v.VehicleHistoryRepo.GetAllVehicleHistory(tx, filterCondition, pages)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (v *VehicleHistoryServiceImpl) GetVehicleHistoryById(id int) (transactionworkshoppayloads.VehicleHistoryByIdResponses, *exceptions.BaseErrorResponse) {
	//TODO implement me
	tx := v.DB.Begin()
	result, err := v.VehicleHistoryRepo.GetVehicleHistoryById(tx, id)
	defer helper.CommitOrRollbackTrx(tx)

	if err != nil {
		return result, err
	}
	return result, nil
}
