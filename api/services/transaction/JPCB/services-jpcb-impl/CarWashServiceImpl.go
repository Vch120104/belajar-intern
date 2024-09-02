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

type CarWashServiceImpl struct {
	CarWashRepository transactionjpcbrepository.CarWashRepository
	DB                *gorm.DB
	RedisClient       *redis.Client
}

func NewCarWashServiceImpl(CarWashRepository transactionjpcbrepository.CarWashRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.CarWashService {
	return &CarWashServiceImpl{
		CarWashRepository: CarWashRepository,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *CarWashServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.CarWashRepository.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *CarWashServiceImpl) UpdatePriority(workOrderSystemNumber int, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.UpdatePriority(tx, workOrderSystemNumber, carWashPriorityId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionjpcbentities.CarWash{}, err
	}
	return result, nil
}

func (s *CarWashServiceImpl) GetAllCarWashPriorityDropDown() ([]transactionjpcbpayloads.CarWashPriorityDropDownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.CarWashRepository.GetAllCarWashPriorityDropDown(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *CarWashServiceImpl) DeleteCarWash(workOrderSystemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.DeleteCarWash(tx, workOrderSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *CarWashServiceImpl) PostCarWash(workOrderSystemNumber int) (transactionjpcbpayloads.CarWashPostResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.PostCarWash(tx, workOrderSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionjpcbpayloads.CarWashPostResponse{}, err
	}
	return result, nil
}

func (s *CarWashServiceImpl) GetAllCarWashScreen(companyId int) ([]transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.CarWashRepository.GetAllCarWashScreen(tx, companyId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *CarWashServiceImpl) UpdateBayNumberCarWashScreen(bayNumber, workOrderSystemNumber int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.UpdateBayNumberCarWashScreen(tx, bayNumber, workOrderSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, err
	}
	return result, nil
}

func (s *CarWashServiceImpl) StartCarWash(workOrderSystemNumber int, carWashBayId int) (transactionjpcbpayloads.CarWashScreenGetAllResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	result, err := s.CarWashRepository.StartCarWash(tx, workOrderSystemNumber, carWashBayId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return transactionjpcbpayloads.CarWashScreenGetAllResponse{}, err
	}
	return result, nil
}
