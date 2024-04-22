package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"

	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemPackageDetailServiceImpl struct {
	ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func StartItemPackageDetailService(ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemPackageDetailService {
	return &ItemPackageDetailServiceImpl{
		ItemPackageDetailRepo: ItemPackageDetailRepo,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailByItemPackageId(tx, itemPackageId, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
