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
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type ClaimSupplierServiceImpl struct {
	claimRepository transactionsparepartrepository.ClaimSupplierRepository
	DB              *gorm.DB
	rdb             *redis.Client
}

func NewClaimSupplierServiceImpl(repo transactionsparepartrepository.ClaimSupplierRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.ClaimSupplierService {
	return &ClaimSupplierServiceImpl{
		claimRepository: repo,
		DB:              db,
		rdb:             rdb,
	}
}
func (service *ClaimSupplierServiceImpl) InsertItemClaim(payload transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
	//tx.Rollback()

	result, err := service.claimRepository.InsertItemClaim(tx, payload)

	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *ClaimSupplierServiceImpl) InsertItemClaimDetail(payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
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
	result, err := service.claimRepository.InsertItemClaimDetail(tx, payloads)
	//tx.Rollback()

	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *ClaimSupplierServiceImpl) GetItemClaimById(itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
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
	result, err := service.claimRepository.GetItemClaimById(tx, itemClaimId)

	if err != nil {
		return result, err
	}
	return result, nil
}
func (service *ClaimSupplierServiceImpl) GetItemClaimDetailByHeaderId(Paginations pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
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
	result, err := service.claimRepository.GetItemClaimDetailByHeaderId(tx, Paginations, filter)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *ClaimSupplierServiceImpl) SubmitItemClaim(claimId int) (bool, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
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
	result, err := service.claimRepository.SubmitItemClaim(tx, claimId)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (service *ClaimSupplierServiceImpl) GetAllItemClaim(page pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := service.DB.Begin()
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
	result, err := service.claimRepository.GetAllItemClaim(tx, page, filter)

	if err != nil {
		return result, err
	}
	return result, nil
}
