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
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
func (s *StockTransactionReasonServiceImpl) GetStockTransactionReasonByCode(Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	result, err := s.StockTransactionRepo.GetStockTransactionReasonByCode(tx, Code)
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
