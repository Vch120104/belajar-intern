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

type ServiceWorkshopServiceImpl struct {
	ServiceWorkshopRepository transactionworkshoprepository.ServiceWorkshopRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func OpenServiceWorkshopServiceImpl(ServiceWorkshopRepo transactionworkshoprepository.ServiceWorkshopRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceWorkshopService {
	return &ServiceWorkshopServiceImpl{
		ServiceWorkshopRepository: ServiceWorkshopRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *ServiceWorkshopServiceImpl) GetAllByTechnicianWO(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ServiceWorkshopDetailResponse, *exceptions.BaseErrorResponse) {

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

	results, repoErr := s.ServiceWorkshopRepository.GetAllByTechnicianWO(tx, idTech, idSysWo, filterCondition, pages)
	if repoErr != nil {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, repoErr
	}

	return results, nil
}

func (s *ServiceWorkshopServiceImpl) StartService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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

	// Start the service
	start, err := s.ServiceWorkshopRepository.StartService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return start, nil
}

func (s *ServiceWorkshopServiceImpl) PendingService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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

	// Pending the service
	pending, err := s.ServiceWorkshopRepository.PendingService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return pending, nil
}

func (s *ServiceWorkshopServiceImpl) TransferService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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

	// Transfer the service
	transfer, err := s.ServiceWorkshopRepository.TransferService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return transfer, nil
}

func (s *ServiceWorkshopServiceImpl) StopService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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

	// Stop the service
	stop, err := s.ServiceWorkshopRepository.StopService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return stop, nil
}
