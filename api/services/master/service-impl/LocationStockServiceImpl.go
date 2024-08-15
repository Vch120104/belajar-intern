package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LocationStockServiceImpl struct {
	LocationStockRepo masterrepository.LocationStockRepository
	DB                *gorm.DB
	RedisClient       *redis.Client
}

func (l *LocationStockServiceImpl) GetAllLocationStock(conditions []utils.FilterCondition, pagination pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := l.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := l.LocationStockRepo.GetAllStock(tx, conditions, pagination)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}

func NewLocationStockServiceImpl(LocationStockService masterrepository.LocationStockRepository, db *gorm.DB, redis *redis.Client) masterservice.LocationStockService {
	return &LocationStockServiceImpl{
		LocationStockRepo: LocationStockService,
		DB:                db,
		RedisClient:       redis,
	}
}
