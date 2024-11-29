package masterserviceimpl

import (
	// "after-sales/api/exceptions"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IncentiveGroupDetailImpl struct {
	IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository
	DB                             *gorm.DB
	RedisClient                    *redis.Client // Redis client
}

func StartIncentiveGroupDetailService(IncentiveGroupDetailRepository masterrepository.IncentiveGroupDetailRepository, db *gorm.DB, redisClient *redis.Client) masterservice.IncentiveGroupDetailService {
	return &IncentiveGroupDetailImpl{
		IncentiveGroupDetailRepository: IncentiveGroupDetailRepository,
		DB:                             db,
		RedisClient:                    redisClient,
	}
}

func (s *IncentiveGroupDetailImpl) GetAllIncentiveGroupDetail(headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.IncentiveGroupDetailRepository.GetAllIncentiveGroupDetail(tx, headerId, pages)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupDetailImpl) GetIncentiveGroupDetailById(id int) (masterpayloads.IncentiveGroupDetailResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupDetailImpl) SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse) {
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

	if req.IncentiveGroupDetailId != 0 {
		_, err := s.IncentiveGroupDetailRepository.GetIncentiveGroupDetailById(tx, req.IncentiveGroupDetailId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.IncentiveGroupDetailRepository.SaveIncentiveGroupDetail(tx, req)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveGroupDetailImpl) UpdateIncentiveGroupDetail(id int, req masterpayloads.UpdateIncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse) {
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

	results, err := s.IncentiveGroupDetailRepository.UpdateIncentiveGroupDetail(tx, id, req)

	if err != nil {
		return results, err
	}
	return results, nil
}
