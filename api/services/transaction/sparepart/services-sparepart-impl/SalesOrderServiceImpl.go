package transactionsparepartserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SalesOrderServiceImpl struct {
	salesOrderRepo transactionsparepartrepository.SalesOrderRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartSalesOrderService(salesOrderRepo transactionsparepartrepository.SalesOrderRepository, db *gorm.DB, redisClient *redis.Client) transactionsparepartservice.SalesOrderService {
	return &SalesOrderServiceImpl{
		salesOrderRepo: salesOrderRepo,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *SalesOrderServiceImpl) GetSalesOrderByID(tx *gorm.DB, id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse) {
	value, err := s.salesOrderRepo.GetSalesOrderByID(tx, id)
	if err != nil {
		return transactionsparepartpayloads.SalesOrderResponse{}, err
	}
	defer helper.CommitOrRollback(tx, err)
	return value, nil
}
