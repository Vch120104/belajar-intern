package transactionjpcbserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.OutstandingJobAllocationRepository.GetAllOutstandingJobAllocation(tx, filterCondition, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *OutstandingJobAllocationServiceImpl) GetByTypeIdOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int) (transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.OutstandingJobAllocationRepository.GetByTypeIdOutstandingJobAllocation(tx, referenceDocumentType, referenceSystemNumber)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *OutstandingJobAllocationServiceImpl) SaveOutstandingJobAllocation(referenceDocumentType string, referenceSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationSaveRequest) (transactionworkshopentities.WorkOrderAllocation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result := transactionworkshopentities.WorkOrderAllocation{}

	operationCodeResult, err := s.OptionCodeRepository.GetOperationCodeById(tx, req.OperationId)
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
