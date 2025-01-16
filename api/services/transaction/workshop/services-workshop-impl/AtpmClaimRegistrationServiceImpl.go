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

type AtpmClaimRegistrationServiceImpl struct {
	AtpmClaimRegistrationRepository transactionworkshoprepository.AtpmClaimRegistrationRepository
	Db                              *gorm.DB
	RedisClient                     *redis.Client // Redis client
}

func OpenAtpmClaimRegistrationServiceImpl(AtpmClaimRegistrationRepository transactionworkshoprepository.AtpmClaimRegistrationRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.AtpmClaimRegistrationService {
	return &AtpmClaimRegistrationServiceImpl{
		AtpmClaimRegistrationRepository: AtpmClaimRegistrationRepository,
		Db:                              db,
		RedisClient:                     redisClient,
	}
}

func (s *AtpmClaimRegistrationServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

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

	results, repoErr := s.AtpmClaimRegistrationRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *AtpmClaimRegistrationServiceImpl) GetById(id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse) {
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

	result, repoErr := s.AtpmClaimRegistrationRepository.GetById(tx, id, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *AtpmClaimRegistrationServiceImpl) New(request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse) {
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

	result, repoErr := s.AtpmClaimRegistrationRepository.New(tx, request)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *AtpmClaimRegistrationServiceImpl) Save(id int, request transactionworkshoppayloads.AtpmClaimRegistrationRequestSave) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse) {
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

	result, repoErr := s.AtpmClaimRegistrationRepository.Save(tx, id, request)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *AtpmClaimRegistrationServiceImpl) Submit(id int) (bool, *exceptions.BaseErrorResponse) {
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

	_, repoErr := s.AtpmClaimRegistrationRepository.Submit(tx, id)
	if repoErr != nil {
		return false, repoErr
	}

	return true, nil
}

func (s *AtpmClaimRegistrationServiceImpl) Void(id int) (bool, *exceptions.BaseErrorResponse) {
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

	_, repoErr := s.AtpmClaimRegistrationRepository.Void(tx, id)
	if repoErr != nil {
		return false, repoErr
	}

	return true, nil
}

func (s *AtpmClaimRegistrationServiceImpl) GetAllServiceHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, repoErr := s.AtpmClaimRegistrationRepository.GetAllServiceHistory(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}
