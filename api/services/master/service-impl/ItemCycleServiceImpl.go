package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemCycleServiceImpl struct {
	ItemCycleRepo masterrepository.ItemCycleRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func NewItemCycleServiceImpl(ItemCycleRepo masterrepository.ItemCycleRepository, db *gorm.DB, rdb *redis.Client) masterservice.ItemCycleService {
	return &ItemCycleServiceImpl{
		ItemCycleRepo: ItemCycleRepo,
		DB:            db,
		RedisClient:   rdb,
	}
}
func (i *ItemCycleServiceImpl) ItemCycleInsert(payloads masterpayloads.ItemCycleInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
	results, err := i.ItemCycleRepo.InsertItemCycle(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return results, err
	}
	return results, nil
}
