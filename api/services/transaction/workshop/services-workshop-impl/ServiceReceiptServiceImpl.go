package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
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

type ServiceReceiptServiceImpl struct {
	ServiceReceiptRepository transactionworkshoprepository.ServiceReceiptRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenServiceReceiptServiceImpl(ServiceReceiptRepo transactionworkshoprepository.ServiceReceiptRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceReceiptService {
	return &ServiceReceiptServiceImpl{
		ServiceReceiptRepository: ServiceReceiptRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *ServiceReceiptServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

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

	results, repoErr := s.ServiceReceiptRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *ServiceReceiptServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			logrus.Info("Transaction rollback due to error:", errResponse)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				errResponse = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	result, repoErr := s.ServiceReceiptRepository.GetById(tx, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *ServiceReceiptServiceImpl) Save(id int, request transactionworkshoppayloads.ServiceReceiptSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
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
			logrus.Info("Transaction rollback due to error:", errResponse)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				errResponse = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	save, repoErr := s.ServiceReceiptRepository.Save(tx, id, request)
	if repoErr != nil {
		return transactionworkshopentities.ServiceRequest{}, repoErr
	}

	return save, nil
}
