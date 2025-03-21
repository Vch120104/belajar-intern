package transactionbodyshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoprepository "after-sales/api/repositories/transaction/bodyshop"
	transactionbodyshopservice "after-sales/api/services/transaction/bodyshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceBodyshopServiceImpl struct {
	ServiceBodyshopRepository transactionbodyshoprepository.ServiceBodyshopRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func OpenServiceBodyshopServiceImpl(ServiceBodyshopRepo transactionbodyshoprepository.ServiceBodyshopRepository, db *gorm.DB, redisClient *redis.Client) transactionbodyshopservice.ServiceBodyshopService {
	return &ServiceBodyshopServiceImpl{
		ServiceBodyshopRepository: ServiceBodyshopRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *ServiceBodyshopServiceImpl) GetAllByTechnicianWOBodyshop(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

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

	results, repoErr := s.ServiceBodyshopRepository.GetAllByTechnicianWOBodyshop(tx, idTech, idSysWo, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func (s *ServiceBodyshopServiceImpl) StartService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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
	start, err := s.ServiceBodyshopRepository.StartService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return start, nil
}

func (s *ServiceBodyshopServiceImpl) PendingService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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
	pending, err := s.ServiceBodyshopRepository.PendingService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return pending, nil
}

func (s *ServiceBodyshopServiceImpl) TransferService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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
	transfer, err := s.ServiceBodyshopRepository.TransferService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return transfer, nil
}

func (s *ServiceBodyshopServiceImpl) StopService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
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
	stop, err := s.ServiceBodyshopRepository.StopService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return stop, nil
}
