package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StockTransactionServiceImpl struct {
	StockTransactionRepo masterrepository.StockTransactionTypeRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func NewStockTransactionServiceImpl(StockTransactionRepo masterrepository.StockTransactionTypeRepository, DB *gorm.DB, redisClient *redis.Client) masterservice.StockTransactionTypeService {
	return &StockTransactionServiceImpl{StockTransactionRepo: StockTransactionRepo, DB: DB, RedisClient: redisClient}
}
func (service *StockTransactionServiceImpl) GetStockTransactionTypeByCode(Code string) (masterentities.StockTransactionType, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.GetStockTransactionTypeByCode(tx, Code)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *StockTransactionServiceImpl) GetAllStockTransactionType(conditions []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.GetAllStockTransactionType(tx, conditions, pagination)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
