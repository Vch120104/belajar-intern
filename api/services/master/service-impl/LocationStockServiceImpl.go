package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
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

type LocationStockServiceImpl struct {
	LocationStockRepo masterrepository.LocationStockRepository
	DB                *gorm.DB
	RedisClient       *redis.Client
}

func (s *LocationStockServiceImpl) UpdateLocationStock(payloads masterwarehousepayloads.LocationStockUpdatePayloads) (bool, *exceptions.BaseErrorResponse) {
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
	results, repoErr := s.LocationStockRepo.UpdateLocationStock(tx, payloads)
	if repoErr != nil {
		return results, repoErr
	}
	return results, nil
}

func (l *LocationStockServiceImpl) GetAllLocationStock(conditions []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := l.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := l.LocationStockRepo.GetAllStock(tx, conditions, pagination)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func NewLocationStockServiceImpl(LocationStockService masterrepository.LocationStockRepository, db *gorm.DB, redis *redis.Client) masterservice.LocationStockService {
	return &LocationStockServiceImpl{
		LocationStockRepo: LocationStockService,
		DB:                db,
		RedisClient:       redis,
	}
}
