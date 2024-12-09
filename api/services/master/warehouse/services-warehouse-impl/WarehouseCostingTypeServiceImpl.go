package masterwarehouseserviceimpl

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WarehouseCostingTypeServiceImpl struct {
	costingTypeRepo masterwarehouserepository.WarehouseCostingTypeRepository
	DB              *gorm.DB
	RedisClient     *redis.Client // Redis client
}

func NewWarehouseCostingTypeServiceImpl(costingTypeRepo masterwarehouserepository.WarehouseCostingTypeRepository, DB *gorm.DB, RedisClient *redis.Client) masterwarehouseservice.WarehouseCostingTypeService {
	return &WarehouseCostingTypeServiceImpl{costingTypeRepo: costingTypeRepo, DB: DB, RedisClient: RedisClient}
}

func (s *WarehouseCostingTypeServiceImpl) GetByCodeWarehouseCostingType(CostingCode string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse) {
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
	get, err := s.costingTypeRepo.GetByCodeWarehouseCostingType(tx, CostingCode)
	if err != nil {
		return get, err
	}
	return get, nil
}
