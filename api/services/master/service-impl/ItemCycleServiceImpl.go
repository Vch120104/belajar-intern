package masterserviceimpl

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemCycleServiceImpl struct {
	ItemCycleRepo masterrepository.ItemCycleRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func NewItemCycleServiceImpl(ItemCycleRepo masterrepository.ItemCycleRepository, db *gorm.DB, rdb *redis.Client) masterservice.ItemCycleService {
	return &ItemCycleServiceImpl{
		ItemCycleRepo: ItemCycleRepo,
		DB:            db,
		RedisClient:   rdb,
	}
}
func (s *ItemCycleServiceImpl) ItemCycleInsert(payloads masterpayloads.ItemCycleInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
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
	results, err := s.ItemCycleRepo.InsertItemCycle(tx, payloads)
	if err != nil {
		return results, err
	}
	return results, nil
}
