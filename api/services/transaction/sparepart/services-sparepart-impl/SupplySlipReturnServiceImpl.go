package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

	results, err := s.supplySlipReturnRepo.SaveSupplySlipReturn(tx, req)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturn{}, err
	}
	return results, nil
}

func (s *SupplySlipReturnServiceImpl) SaveSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse) {
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

	results, err := s.supplySlipReturnRepo.SaveSupplySlipReturnDetail(tx, req)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturnDetail{}, err
	}
	return results, nil
}

func (s *SupplySlipReturnServiceImpl) GetAllSupplySlipReturn(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
	results, totalPages, totalRows, err := s.supplySlipReturnRepo.GetAllSupplySlipReturn(tx, internalFilter, externalFilter, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *SupplySlipReturnServiceImpl) GetSupplySlipReturnById(Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	// Memulai transaksi
	tx := s.DB.Begin()
	var errResponse *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errResponse = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if errResponse != nil {
			tx.Rollback()
			logrus.WithError(errResponse.Err).Info("Transaction rollback due to error")
		} else {
			if err := tx.Commit().Error; err != nil {
				logrus.WithError(err).Error("Transaction commit failed")
				errResponse = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", err),
				}
			}
		}
	}()

	supplySlipId, err := s.supplySlipReturnRepo.GetSupplySlipId(tx, Id)
	if err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        fmt.Errorf("failed to fetch supply slip ID: %w", err),
		}
		return nil, errResponse
	}

	supplyResults, err := s.supplySlipRepo.GetSupplySlipById(tx, supplySlipId, pagination)
	if err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch supply slip: %w", err),
		}
		return nil, errResponse
	}

	results, err := s.supplySlipReturnRepo.GetSupplySlipReturnById(tx, Id, pagination, supplyResults)
	if err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch supply slip return: %w", err),
		}
		return nil, errResponse
	}

	return results, nil
}

func (s *SupplySlipReturnServiceImpl) GetSupplySlipReturnDetailById(id int) (transactionsparepartpayloads.SupplySlipReturnDetailResponse, *exceptions.BaseErrorResponse) {
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
	value, err := s.supplySlipReturnRepo.GetSupplySlipReturnDetailById(tx, id)
	if err != nil {
		return transactionsparepartpayloads.SupplySlipReturnDetailResponse{}, err
	}
	return value, nil
}

func (s *SupplySlipReturnServiceImpl) UpdateSupplySlipReturn(req transactionsparepartentities.SupplySlipReturn, id int) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse) {
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
	result, err := s.supplySlipReturnRepo.UpdateSupplySlipReturn(tx, req, id)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturn{}, err
	}

	return result, nil
}

func (s *SupplySlipReturnServiceImpl) UpdateSupplySlipReturnDetail(req transactionsparepartentities.SupplySlipReturnDetail, id int) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse) {
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
	result, err := s.supplySlipReturnRepo.UpdateSupplySlipReturnDetail(tx, req, id)
	if err != nil {
		return transactionsparepartentities.SupplySlipReturnDetail{}, err
	}

	return result, nil
}
