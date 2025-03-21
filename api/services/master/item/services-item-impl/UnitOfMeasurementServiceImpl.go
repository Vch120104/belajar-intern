package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UnitOfMeasurementServiceImpl struct {
	unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func StartUnitOfMeasurementService(unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.UnitOfMeasurementService {
	return &UnitOfMeasurementServiceImpl{
		unitOfMeasurementRepo: unitOfMeasurementRepo,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurementIsActive(tx)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementById(id int) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementByCode(tx, Code)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurement(tx, filterCondition, pages)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) ChangeStatusUnitOfMeasurement(Id int) (bool, *exceptions.BaseErrorResponse) {
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
	_, err = s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.unitOfMeasurementRepo.ChangeStatusUnitOfMeasurement(tx, Id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, *exceptions.BaseErrorResponse) {
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

	if req.UomId != 0 {
		_, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, req.UomId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.unitOfMeasurementRepo.SaveUnitOfMeasurement(tx, req)

	if err != nil {
		return false, err
	}
	return results, nil
}
func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementItem(payload masteritempayloads.UomItemRequest) (masteritempayloads.UomItemResponses, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementItem(tx, payload)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetQuantityConversion(payload masteritempayloads.UomGetQuantityConversion) (masteritempayloads.GetQuantityConversionResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.unitOfMeasurementRepo.GetQuantityConversion(tx, payload)

	if err != nil {
		return results, err
	}
	return results, nil
}
