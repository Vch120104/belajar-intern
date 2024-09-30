package transactionjpcbserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OutstandingJobAllocationServiceImpl struct {
	OutstandingJobAllocationRepository transactionjpcbrepository.OutstandingJobAllocationRepository
	OptionCodeRepository               masteroperationrepository.OperationCodeRepository
	DB                                 *gorm.DB
	Redis                              *redis.Client
}

func StartOutstandingJobAllocationService(outstandingJobAllocationRepository transactionjpcbrepository.OutstandingJobAllocationRepository, optionCodeRepository masteroperationrepository.OperationCodeRepository, db *gorm.DB, redis *redis.Client) transactionjpcbservice.OutstandingJobAllocationService {
	return &OutstandingJobAllocationServiceImpl{
		OutstandingJobAllocationRepository: outstandingJobAllocationRepository,
		OptionCodeRepository:               optionCodeRepository,
		DB:                                 db,
		Redis:                              redis,
	}
}

func (s *OutstandingJobAllocationServiceImpl) GetAllOutstandingJobAllocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.OutstandingJobAllocationRepository.GetAllOutstandingJobAllocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *OutstandingJobAllocationServiceImpl) GetByTypeIdOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int) (transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.OutstandingJobAllocationRepository.GetByTypeIdOutstandingJobAllocation(tx, referenceDocumentType, referenceSystemNumber)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *OutstandingJobAllocationServiceImpl) SaveOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationSaveRequest) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result := transactionjpcbpayloads.SettingTechnicianGetByIdResponse{}

	operationCodeResult, err := s.OptionCodeRepository.GetOperationCodeById(tx, req.OperationId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}

	result, updateRequest, err := s.OutstandingJobAllocationRepository.SaveOutstandingJobAllocation(tx, referenceDocumentType, referenceSystemNumber, req, operationCodeResult)
	if err != nil {
		return result, err
	}

	if updateRequest.TechAllocSystemNumber != 0 {
		updateResult, err := s.OutstandingJobAllocationRepository.UpdateOutstandingJobAllocation(tx, updateRequest.TechAllocSystemNumber, updateRequest)
		if err != nil {
			return result, err
		}

		recalculateErr := s.OutstandingJobAllocationRepository.ReCalculateTimeJob(tx, updateResult.SourceFirstTechAllocSystenNumber)
		if recalculateErr != nil {
			return result, err
		}

		recalculateErr = s.OutstandingJobAllocationRepository.ReCalculateTimeJob(tx, updateResult.TechAllocSystemNumber)
		if recalculateErr != nil {
			return result, err
		}
	}

	return result, nil
}
