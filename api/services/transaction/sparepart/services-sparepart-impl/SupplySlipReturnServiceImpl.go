package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SupplySlipReturnServiceImpl struct {
	supplySlipReturnRepo transactionsparepartrepository.SupplySlipReturnRepository
	supplySlipRepo       transactionsparepartrepository.SupplySlipRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func StartSupplySlipReturnService(supplySlipReturnRepo transactionsparepartrepository.SupplySlipReturnRepository, supplySlipRepo transactionsparepartrepository.SupplySlipRepository, db *gorm.DB, redisClient *redis.Client) transactionsparepartservice.SupplySlipReturnService {
	return &SupplySlipReturnServiceImpl{
		supplySlipReturnRepo: supplySlipReturnRepo,
		supplySlipRepo:       supplySlipRepo,
		DB:                   db,
		RedisClient:          redisClient,
	}
}

func (s *SupplySlipReturnServiceImpl) SaveSupplySlipReturn(req transactionsparepartentities.SupplySlipReturn) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.supplySlipReturnRepo.SaveSupplySlipReturn(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturn{}, err
	}
	return results, nil
}

func (s *SupplySlipReturnServiceImpl) SaveSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.supplySlipReturnRepo.SaveSupplySlipReturnDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturnDetail{}, err
	}
	return results, nil
}

func (s *SupplySlipReturnServiceImpl) GetAllSupplySlipReturn(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.supplySlipReturnRepo.GetAllSupplySlipReturn(tx, internalFilter, externalFilter, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *SupplySlipReturnServiceImpl) GetSupplySlipReturnById(Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	supplySlipId, errSupplyId := s.supplySlipReturnRepo.GetSupplySlipId(tx, Id)

	if errSupplyId != nil {
		return nil, errSupplyId
	} 

	supplyResults, errSupply := s.supplySlipRepo.GetSupplySlipById(tx, supplySlipId, pagination)

	if errSupply != nil {
		return nil, errSupply
	} 

	results, err := s.supplySlipReturnRepo.GetSupplySlipReturnById(tx, Id, pagination, supplyResults)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SupplySlipReturnServiceImpl) GetSupplySlipReturnDetailById(id int) (transactionsparepartpayloads.SupplySlipReturnDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	value, err := s.supplySlipReturnRepo.GetSupplySlipReturnDetailById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipReturnDetailResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipReturnServiceImpl) UpdateSupplySlipReturn(req transactionsparepartentities.SupplySlipReturn, id int)(transactionsparepartentities.SupplySlipReturn,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.supplySlipReturnRepo.UpdateSupplySlipReturn(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return transactionsparepartentities.SupplySlipReturn{},err
	}

	return result, nil
}

func (s *SupplySlipReturnServiceImpl) UpdateSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail, id int)(transactionsparepartentities.SupplySlipReturnDetail,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.supplySlipReturnRepo.UpdateSupplySlipReturnDetail(tx,req,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return transactionsparepartentities.SupplySlipReturnDetail{},err
	}

	return result, nil
}