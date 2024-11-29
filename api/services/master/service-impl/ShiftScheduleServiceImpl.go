package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShiftScheduleServiceImpl struct {
	ShiftScheduleRepo masterrepository.ShiftScheduleRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartShiftScheduleService(ShiftScheduleRepo masterrepository.ShiftScheduleRepository, db *gorm.DB, redisClient *redis.Client) masterservice.ShiftScheduleService {
	return &ShiftScheduleServiceImpl{
		ShiftScheduleRepo: ShiftScheduleRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *ShiftScheduleServiceImpl) GetShiftScheduleById(id int) (masterpayloads.ShiftScheduleResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.ShiftScheduleRepo.GetAllShiftSchedule(tx, filterCondition, pages)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) ChangeStatusShiftSchedule(oprId int) (bool, *exceptions.BaseErrorResponse) {
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

	_, err = s.ShiftScheduleRepo.GetShiftScheduleById(tx, oprId)

	if err != nil {
		return false, err
	}

	results, err := s.ShiftScheduleRepo.ChangeStatusShiftSchedule(tx, oprId)

	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) SaveShiftSchedule(req masterpayloads.ShiftScheduleResponse) (bool, *exceptions.BaseErrorResponse) {
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

	if req.ShiftScheduleId != 0 {
		_, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, req.ShiftScheduleId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.ShiftScheduleRepo.SaveShiftSchedule(tx, req)

	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) GetShiftScheduleDropDown() ([]masterpayloads.ShiftScheduleDropDownResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.ShiftScheduleRepo.GetShiftScheduleDropDown(tx)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) UpdateShiftSchedule(Id int, request masterpayloads.ShiftScheduleUpdate) (masterentities.ShiftSchedule, *exceptions.BaseErrorResponse) {
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
	results, err := s.ShiftScheduleRepo.UpdateShiftSchedule(tx, Id, request)

	if err != nil {
		return masterentities.ShiftSchedule{}, err
	}
	return results, nil
}
