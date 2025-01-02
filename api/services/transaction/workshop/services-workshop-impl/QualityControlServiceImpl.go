package transactionworkshopserviceimpl

import (
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

type QualityControlServiceImpl struct {
	QualityControlRepository transactionworkshoprepository.QualityControlRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenQualityControlServiceImpl(QualityControlRepo transactionworkshoprepository.QualityControlRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.QualityControlService {
	return &QualityControlServiceImpl{
		QualityControlRepository: QualityControlRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *QualityControlServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, repoErr := s.QualityControlRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *QualityControlServiceImpl) GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.GetById(tx, id, filterCondition, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlServiceImpl) Qcpass(id int, iddet int) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.Qcpass(tx, id, iddet)
	if repoErr != nil {

		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlServiceImpl) Reorder(id int, iddet int, payload transactionworkshoppayloads.QualityControlReorder) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.Reorder(tx, id, iddet, payload)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}
