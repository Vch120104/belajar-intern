package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationDefinitionServiceImpl struct {
	WarehouseLocationDefinitionRepo masterwarehouserepository.WarehouseLocationDefinitionRepository
	DB                              *gorm.DB
	RedisClient                     *redis.Client // Redis client
}

func OpenWarehouseLocationDefinitionService(WarehouseLocationDefinition masterwarehouserepository.WarehouseLocationDefinitionRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseLocationDefinitionService {
	return &WarehouseLocationDefinitionServiceImpl{
		WarehouseLocationDefinitionRepo: WarehouseLocationDefinition,
		DB:                              db,
		RedisClient:                     redisClient,
	}
}

func (s *WarehouseLocationDefinitionServiceImpl) Save(request masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse) {
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

	if request.WarehouseLocationDefinitionId != 0 {
		_, err := s.WarehouseLocationDefinitionRepo.GetById(tx, request.WarehouseLocationDefinitionId)

		if err != nil {
			return masterwarehouseentities.WarehouseLocationDefinition{}, err
		}
	}

	save, err := s.WarehouseLocationDefinitionRepo.Save(tx, request)

	if err != nil {
		return masterwarehouseentities.WarehouseLocationDefinition{}, err
	}
	return save, err
}

func (s *WarehouseLocationDefinitionServiceImpl) SaveData(request masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse) {
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

	if request.WarehouseLocationDefinitionId != 0 {
		_, err := s.WarehouseLocationDefinitionRepo.GetById(tx, request.WarehouseLocationDefinitionId)

		if err != nil {
			return masterwarehouseentities.WarehouseLocationDefinition{}, err
		}
	}

	save, err := s.WarehouseLocationDefinitionRepo.SaveData(tx, request)

	if err != nil {
		return masterwarehouseentities.WarehouseLocationDefinition{}, err
	}
	return save, err
}

func (s *WarehouseLocationDefinitionServiceImpl) GetById(Id int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.WarehouseLocationDefinitionRepo.GetById(tx, Id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) GetByLevel(idlevel int, idwhl string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.WarehouseLocationDefinitionRepo.GetByLevel(tx, idlevel, idwhl)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
	results, totalPages, totalRows, err := s.WarehouseLocationDefinitionRepo.GetAll(tx, filterCondition, pages)

	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) ChangeStatus(Id int) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse) {
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

	// Ubah status
	entity, err := s.WarehouseLocationDefinitionRepo.ChangeStatus(tx, Id)

	if err != nil {
		return masterwarehouseentities.WarehouseLocationDefinition{}, err
	}
	return entity, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) PopupWarehouseLocationLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
	results, totalPages, totalRows, err := s.WarehouseLocationDefinitionRepo.PopupWarehouseLocationLevel(tx, filterCondition, pages)

	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}
