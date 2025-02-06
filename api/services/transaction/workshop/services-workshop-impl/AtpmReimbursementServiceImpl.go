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

type AtpmReimbursementServiceImpl struct {
	AtpmReimbursementRepository transactionworkshoprepository.AtpmReimbursementRepository
	Db                          *gorm.DB
	RedisClient                 *redis.Client // Redis client
}

func OpenAtpmReimbursementServiceImpl(AtpmReimbursementRepository transactionworkshoprepository.AtpmReimbursementRepository, Db *gorm.DB, RedisClient *redis.Client) transactionworkshopservice.AtpmReimbursementService {
	return &AtpmReimbursementServiceImpl{
		AtpmReimbursementRepository: AtpmReimbursementRepository,
		Db:                          Db,
		RedisClient:                 RedisClient,
	}
}

func (s *AtpmReimbursementServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
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

	results, repoErr := s.AtpmReimbursementRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *AtpmReimbursementServiceImpl) New(req transactionworkshoppayloads.AtpmReimbursementRequest) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
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

	result, repoErr := s.AtpmReimbursementRepository.New(tx, req)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *AtpmReimbursementServiceImpl) Save(claimsysno int, req transactionworkshoppayloads.AtpmReimbursementUpdate) (transactionworkshopentities.AtpmReimbursement, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
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

	result, repoErr := s.AtpmReimbursementRepository.Save(tx, claimsysno, req)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *AtpmReimbursementServiceImpl) Submit(claimsysno int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
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

	result, repoErr := s.AtpmReimbursementRepository.Submit(tx, claimsysno)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}
