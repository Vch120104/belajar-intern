package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
func (s *VehicleHistoryServiceImpl) GetAllVehicleHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//TODO implement me
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
	result, err := s.VehicleHistoryRepo.GetAllVehicleHistory(tx, filterCondition, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *VehicleHistoryServiceImpl) GetVehicleHistoryById(id int) (transactionworkshoppayloads.VehicleHistoryByIdResponses, *exceptions.BaseErrorResponse) {
	//TODO implement me
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
	result, err := s.VehicleHistoryRepo.GetVehicleHistoryById(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}
