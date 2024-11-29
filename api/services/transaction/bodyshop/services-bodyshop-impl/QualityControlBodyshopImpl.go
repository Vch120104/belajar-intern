package transactionbodyshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	transactionbodyshoprepository "after-sales/api/repositories/transaction/bodyshop"
	transactionbodyshopservice "after-sales/api/services/transaction/bodyshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QualityControlBodyshopServiceImpl struct {
	QualityControlRepository transactionbodyshoprepository.QualityControlBodyshopRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenQualityControlBodyshopServiceImpl(QualityControlRepo transactionbodyshoprepository.QualityControlBodyshopRepository, db *gorm.DB, redisClient *redis.Client) transactionbodyshopservice.QualityControlBodyshopService {
	return &QualityControlBodyshopServiceImpl{
		QualityControlRepository: QualityControlRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *QualityControlBodyshopServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

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

	results, totalPages, totalRows, repoErr := s.QualityControlRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (s *QualityControlBodyshopServiceImpl) GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.GetById(tx, id, filterCondition, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlBodyshopServiceImpl) Qcpass(id int, iddet int) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.Qcpass(tx, id, iddet)
	if repoErr != nil {

		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlBodyshopServiceImpl) Reorder(id int, iddet int, payload transactionbodyshoppayloads.QualityControlReorder) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

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

	result, repoErr := s.QualityControlRepository.Reorder(tx, id, iddet, payload)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}
