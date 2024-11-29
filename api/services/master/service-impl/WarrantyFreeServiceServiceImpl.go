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

type WarrantyFreeServiceServiceImpl struct {
	warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository
	DB                      *gorm.DB
	RedisClient             *redis.Client // Redis client
}

func StartWarrantyFreeServiceService(warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository, db *gorm.DB, redisClient *redis.Client) masterservice.WarrantyFreeServiceService {
	return &WarrantyFreeServiceServiceImpl{
		warrantyFreeServiceRepo: warrantyFreeServiceRepo,
		DB:                      db,
		RedisClient:             redisClient,
	}
}

func (s *WarrantyFreeServiceServiceImpl) GetAllWarrantyFreeService(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.warrantyFreeServiceRepo.GetAllWarrantyFreeService(tx, filterCondition, pages)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) GetWarrantyFreeServiceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
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
	results, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse) {
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

	results, err := s.warrantyFreeServiceRepo.SaveWarrantyFreeService(tx, req)

	if err != nil {
		return masterentities.WarrantyFreeService{}, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) ChangeStatusWarrantyFreeService(Id int) (masterpayloads.WarrantyFreeServicePatchResponse, *exceptions.BaseErrorResponse) {
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

	_, err = s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)

	if err != nil {
		return masterpayloads.WarrantyFreeServicePatchResponse{}, err
	}

	results, err := s.warrantyFreeServiceRepo.ChangeStatusWarrantyFreeService(tx, Id)

	if err != nil {
		return masterpayloads.WarrantyFreeServicePatchResponse{}, err
	}
	return results, nil
}

func (s *WarrantyFreeServiceServiceImpl) UpdateWarrantyFreeService(req masterentities.WarrantyFreeService, id int) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse) {
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
	result, err := s.warrantyFreeServiceRepo.UpdateWarrantyFreeService(tx, req, id)

	if err != nil {
		return masterentities.WarrantyFreeService{}, err
	}

	return result, nil
}
