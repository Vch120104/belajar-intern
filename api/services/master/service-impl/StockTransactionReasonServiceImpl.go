package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
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

func (service *StockTransactionReasonServiceImpl) InsertStockTransactionReason(payloads masterpayloads.StockTransactionReasonInsertPayloads) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.InsertStockTransactionReason(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *StockTransactionReasonServiceImpl) GetStockTransactionReasonById(id int) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.GetStockTransactionReasonById(tx, id)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *StockTransactionReasonServiceImpl) GetAllStockTransactionReason(conditions []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	result, err := service.StockTransactionRepo.GetAllStockTransactionReason(tx, conditions, pagination)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
