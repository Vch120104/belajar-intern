package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StockTransactionReasonServiceImpl struct {
	StockTransactionRepo masterrepository.StockTransactionReasonRepository
	DB                   *gorm.DB
	Redis                *redis.Client
}

func StartStockTransactionReasonServiceImpl(StockTransactionRepo masterrepository.StockTransactionReasonRepository, DB *gorm.DB, Redis *redis.Client) masterservice.StockTransactionReasonService {
	return &StockTransactionReasonServiceImpl{StockTransactionRepo: StockTransactionRepo, DB: DB, Redis: Redis}
}
func (service *StockTransactionReasonServiceImpl) GetStockTransactionReasonByCode(Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.GetStockTransactionReasonByCode(tx, Code)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
