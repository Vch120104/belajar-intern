package transactionsparepartserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SupplySlipServiceImpl struct {
	supplySlipRepo transactionsparepartrepository.SupplySlipRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartSupplySlipService(supplySlipRepo transactionsparepartrepository.SupplySlipRepository, db *gorm.DB, redisClient *redis.Client) transactionsparepartservice.SupplySlipService {
	return &SupplySlipServiceImpl{
		supplySlipRepo: supplySlipRepo,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *SupplySlipServiceImpl) GetSupplySlipById(tx *gorm.DB, id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptionsss_test.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	value, err := s.supplySlipRepo.GetSupplySlipById(tx, id)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipServiceImpl) GetSupplySlipDetailById(tx *gorm.DB, id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptionsss_test.BaseErrorResponse) {
	defer helper.CommitOrRollback(tx)
	value, err := s.supplySlipRepo.GetSupplySlipDetailById(tx, id)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipDetailResponse{}, err
	}
	return value, nil
}
