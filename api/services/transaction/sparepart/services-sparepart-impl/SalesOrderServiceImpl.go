package transactionsparepartserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SalesOrderServiceImpl struct {
	salesOrderRepo transactionsparepartrepository.SalesOrderRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartSalesOrderService(salesOrderRepo transactionsparepartrepository.SalesOrderRepository, db *gorm.DB, redisClient *redis.Client) transactionsparepartservice.SalesOrderService {
	return &SalesOrderServiceImpl{
		salesOrderRepo: salesOrderRepo,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *SalesOrderServiceImpl) GetSalesOrderByID(tx *gorm.DB, id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse) {
	tx = s.DB.Begin()
	var result transactionsparepartpayloads.SalesOrderResponse
	var errResponse *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errResponse = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if errResponse != nil {
			tx.Rollback()
			logrus.WithError(errResponse.Err).Info("Transaction rollback due to error")
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				errResponse = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()

	result, err := s.salesOrderRepo.GetSalesOrderByID(tx, id)
	if err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Sales Order not found",
			Err:        err,
		}
		return transactionsparepartpayloads.SalesOrderResponse{}, errResponse
	}

	return result, nil
}
