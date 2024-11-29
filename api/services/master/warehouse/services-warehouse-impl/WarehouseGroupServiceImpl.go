package masterwarehouseserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type WarehouseGroupServiceImpl struct {
	warehouseGroupRepo masterwarehouserepository.WarehouseGroupRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func OpenWarehouseGroupService(warehouseGroup masterwarehouserepository.WarehouseGroupRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseGroupService {
	return &WarehouseGroupServiceImpl{
		warehouseGroupRepo: warehouseGroup,
		DB:                 db,
		RedisClient:        redisClient,
	}
}

// GetbyGroupCode implements masterwarehouseservice.WarehouseGroupService.
func (s *WarehouseGroupServiceImpl) GetbyGroupCode(groupCode string) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse) {
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
	get, err := s.warehouseGroupRepo.GetbyGroupCode(tx, groupCode)

	if err != nil {
		return get, err
	}
	return get, nil
}

// GetWarehouseGroupDropdownbyId implements masterwarehouseservice.WarehouseGroupService.
func (s *WarehouseGroupServiceImpl) GetWarehouseGroupDropdownbyId(Id int) (masterwarehousepayloads.GetWarehouseGroupDropdown, *exceptions.BaseErrorResponse) {
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
	get, err := s.warehouseGroupRepo.GetWarehouseGroupDropdownbyId(tx, Id)

	if err != nil {
		return get, err
	}
	return get, nil
}

// GetWarehouseGroupDropdown implements masterwarehouseservice.WarehouseGroupService.
func (s *WarehouseGroupServiceImpl) GetWarehouseGroupDropdown() ([]masterwarehousepayloads.GetWarehouseGroupDropdown, *exceptions.BaseErrorResponse) {
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
	get, err := s.warehouseGroupRepo.GetWarehouseGroupDropdown(tx)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseGroupServiceImpl) SaveWarehouseGroup(request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptions.BaseErrorResponse) {
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

	if request.WarehouseGroupId != 0 {
		_, err := s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, request.WarehouseGroupId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.warehouseGroupRepo.SaveWarehouseGroup(tx, request)

	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *WarehouseGroupServiceImpl) GetByIdWarehouseGroup(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse) {
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
	get, err := s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseGroupServiceImpl) GetAllWarehouseGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	get, err := s.warehouseGroupRepo.GetAllWarehouseGroup(tx, filterCondition, pages)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseGroupServiceImpl) ChangeStatusWarehouseGroup(warehouseGroupId int) (bool, *exceptions.BaseErrorResponse) {
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

	_, err = s.warehouseGroupRepo.GetByIdWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return false, err
	}

	change_status, err := s.warehouseGroupRepo.ChangeStatusWarehouseGroup(tx, warehouseGroupId)

	if err != nil {
		return change_status, err
	}
	return change_status, nil
}
