package masterserviceimpl

import (
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StockTransactionServiceImpl struct {
	StockTransactionRepository transactionsparepartrepository.StockTransactionRepository
	DB                         *gorm.DB
	rdb                        *redis.Client
}

func StartStockTransactionServiceImpl(StockTransactionRepository transactionsparepartrepository.StockTransactionRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.StockTransactionService {
	return &StockTransactionServiceImpl{
		StockTransactionRepository: StockTransactionRepository,
		DB:                         db,
		rdb:                        rdb,
	}
}
func (s *StockTransactionServiceImpl) StockTransactionInsert(payloads transactionsparepartpayloads.StockTransactionInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
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
			}
		}
	}()
	result, err := s.StockTransactionRepository.StockTransactionInsert(tx, payloads)

	if err != nil {
		return result, err
	}
	return result, nil
}
