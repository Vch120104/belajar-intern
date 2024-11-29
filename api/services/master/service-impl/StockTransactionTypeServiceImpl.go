package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StockTransactionTypeServiceImpl struct {
	StockTransactionRepo masterservice.StockTransactionTypeRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func NewStockTransactionTypeServiceImpl(StockTransactionRepo masterservice.StockTransactionTypeRepository, DB *gorm.DB, redisClient *redis.Client) masterservice.StockTransactionTypeService {
	return &StockTransactionTypeServiceImpl{StockTransactionRepo: StockTransactionRepo, DB: DB, RedisClient: redisClient}
}
func (s *StockTransactionTypeServiceImpl) GetStockTransactionTypeByCode(Code string) (masterentities.StockTransactionType, *exceptions.BaseErrorResponse) {
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
	result, err := s.StockTransactionRepo.GetStockTransactionTypeByCode(tx, Code)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *StockTransactionTypeServiceImpl) GetAllStockTransactionType(conditions []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.StockTransactionRepo.GetAllStockTransactionType(tx, conditions, pagination)
	if err != nil {
		return result, err
	}
	return result, nil
}
