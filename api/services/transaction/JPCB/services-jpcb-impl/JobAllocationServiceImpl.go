package transactionjpcbserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
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

type JobAllocationServiceImpl struct {
	JobAllocationRepository transactionjpcbrepository.JobAllocationRepository
	DB                      *gorm.DB
	RedisClient             *redis.Client
}

func StartJobAllocationService(jobAllocationRepository transactionjpcbrepository.JobAllocationRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.JobAllocationService {
	return &JobAllocationServiceImpl{
		JobAllocationRepository: jobAllocationRepository,
		DB:                      db,
		RedisClient:             redisClient,
	}
}

func (s *JobAllocationServiceImpl) GetAllJobAllocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.JobAllocationRepository.GetAllJobAllocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *JobAllocationServiceImpl) GetJobAllocationById(technicianAllocationSystemNumber int) (transactionjpcbpayloads.GetJobAllocationByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.JobAllocationRepository.GetJobAllocationById(tx, technicianAllocationSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *JobAllocationServiceImpl) UpdateJobAllocation(technicianAllocationSystemNumber int, req transactionjpcbpayloads.JobAllocationUpdateRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	_, err := s.JobAllocationRepository.GetJobAllocationById(tx, technicianAllocationSystemNumber)
	if err != nil {
		return transactionworkshopentities.WorkOrderAllocation{}, err
	}
	defer helper.CommitOrRollback(tx, err)
	update, err := s.JobAllocationRepository.UpdateJobAllocation(tx, technicianAllocationSystemNumber, req)
	if err != nil {
		return update, err
	}
	return update, nil
}

func (s *JobAllocationServiceImpl) DeleteJobAllocation(technicianAllocationSystemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	_, err := s.JobAllocationRepository.GetJobAllocationById(tx, technicianAllocationSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	delete, err := s.JobAllocationRepository.DeleteJobAllocation(tx, technicianAllocationSystemNumber)
	if err != nil {
		return delete, err
	}
	return delete, nil
}
