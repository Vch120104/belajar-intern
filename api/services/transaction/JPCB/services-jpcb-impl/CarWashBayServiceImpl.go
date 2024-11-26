package transactionjpcbserviceimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BayMasterServiceImpl struct {
	BayMasterRepository transactionjpcbrepository.BayMasterRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client
}

func NewCarWashBayServiceImpl(BayRepository transactionjpcbrepository.BayMasterRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.BayMasterService {
	return &BayMasterServiceImpl{
		BayMasterRepository: BayRepository,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

func (s *BayMasterServiceImpl) GetAllCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.BayMasterRepository.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BayMasterServiceImpl) GetAllActiveCarWashBay(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.BayMasterRepository.GetAllActive(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BayMasterServiceImpl) GetAllDeactiveCarWashBay(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BayMasterRepository.GetAllDeactive(tx, filterCondition)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BayMasterServiceImpl) ChangeStatusCarWashBay(request transactionjpcbpayloads.CarWashBayUpdateRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BayMasterRepository.ChangeStatus(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return transactionjpcbentities.BayMaster{}, err
	}
	return results, nil
}

func (s *BayMasterServiceImpl) GetAllCarWashBayDropDown(filterCondition []utils.FilterCondition) ([]transactionjpcbpayloads.CarWashBayDropDownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BayMasterRepository.CarWashBayDropDown(tx, filterCondition)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *BayMasterServiceImpl) PostCarWashBay(request transactionjpcbpayloads.CarWashBayPostRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.BayMasterRepository.PostCarWashBay(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return transactionjpcbentities.BayMaster{}, err
	}
	return result, nil
}

func (s *BayMasterServiceImpl) PutCarWashBay(request transactionjpcbpayloads.CarWashBayPutRequest) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.BayMasterRepository.UpdateCarWashBay(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return transactionjpcbentities.BayMaster{}, err
	}
	return result, nil
}

func (s *BayMasterServiceImpl) GetCarWashBayById(carWashBayId int) (transactionjpcbentities.BayMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.BayMasterRepository.GetCarWashBayById(tx, carWashBayId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return transactionjpcbentities.BayMaster{}, err
	}
	return result, nil
}
