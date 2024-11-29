package transactionjpcbserviceimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TechnicianAttendanceImpl struct {
	TechnicianAttendanceRepository transactionjpcbrepository.TechnicianAttendanceRepository
	DB                             *gorm.DB
	RedisClient                    *redis.Client
}

func StartTechnicianAttendanceImpl(technicianAttendanceRepository transactionjpcbrepository.TechnicianAttendanceRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.TechnicianAttendanceService {
	return &TechnicianAttendanceImpl{
		TechnicianAttendanceRepository: technicianAttendanceRepository,
		DB:                             db,
		RedisClient:                    redisClient,
	}
}

func (s *TechnicianAttendanceImpl) GetAllTechnicianAttendance(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.TechnicianAttendanceRepository.GetAllTechnicianAttendance(tx, filterCondition, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *TechnicianAttendanceImpl) SaveTechnicianAttendance(req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
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
	result, err := s.TechnicianAttendanceRepository.SaveTechnicianAttendance(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *TechnicianAttendanceImpl) ChangeStatusTechnicianAttendance(technicianAttendanceId int) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
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
	result, err := s.TechnicianAttendanceRepository.ChangeStatusTechnicianAttendance(tx, technicianAttendanceId)
	if err != nil {
		return result, err
	}
	return result, nil
}
