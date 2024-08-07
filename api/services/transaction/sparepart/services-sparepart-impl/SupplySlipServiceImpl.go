package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"

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

func (s *SupplySlipServiceImpl) GetSupplySlipById(tx *gorm.DB, id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse) {
	value, err := s.supplySlipRepo.GetSupplySlipById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipServiceImpl) GetSupplySlipDetailById(tx *gorm.DB, id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse) {
	value, err := s.supplySlipRepo.GetSupplySlipDetailById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipDetailResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipServiceImpl) SaveSupplySlip(req transactionsparepartentities.SupplySlip) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.supplySlipRepo.SaveSupplySlip(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartentities.SupplySlip{}, err
	}
	return results, nil
}

func (s *SupplySlipServiceImpl) SaveSupplySlipDetail(req transactionsparepartentities.SupplySlipDetail) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.supplySlipRepo.SaveSupplySlipDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartentities.SupplySlipDetail{}, err
	}
	return results, nil
}

func (s *SupplySlipServiceImpl) GetAllSupplySlip(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.supplySlipRepo.GetAllSupplySlip(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}