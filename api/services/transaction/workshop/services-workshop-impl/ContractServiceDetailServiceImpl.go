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

type ContractServiceDetailServiceImpl struct {
	ContractServiceDetailRepository transactionworkshoprepository.ContractServiceDetailRepository
	DB                              *gorm.DB
	RedisClient                     *redis.Client
}

func OpenContractServiceDetailServiceImpl(ContractServiceDetailRepo transactionworkshoprepository.ContractServiceDetailRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ContractServiceDetailService {
	return &ContractServiceDetailServiceImpl{
		ContractServiceDetailRepository: ContractServiceDetailRepo,
		DB:                              db,
		RedisClient:                     redisClient,
	}
}

// GetAllDetail implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) GetAllDetail(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	results, repoErr := s.ContractServiceDetailRepository.GetAllDetail(tx, Id, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

// GetById implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) GetById(Id int) (transactionworkshoppayloads.ContractServiceIdResponse, *exceptions.BaseErrorResponse) {
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
	result, err := s.ContractServiceDetailRepository.GetById(tx, Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SaveDetail implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) SaveDetail(req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
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
	result, err := s.ContractServiceDetailRepository.SaveDetail(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateDetail implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) UpdateDetail(contractServiceSystemNumber int, contractServiceLine string, req transactionworkshoppayloads.ContractServiceDetailRequest) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
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
	result, err := s.ContractServiceDetailRepository.UpdateDetail(tx, contractServiceSystemNumber, contractServiceLine, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteDetail implements transactionworkshopservice.ContractServiceDetailService.
func (s *ContractServiceDetailServiceImpl) DeleteDetail(contractServiceSystemNumber int, packageCode string) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.ContractServiceDetailRepository.DeleteDetail(tx, contractServiceSystemNumber, packageCode)
	if err != nil {
		return result, err
	}
	return result, nil
}
