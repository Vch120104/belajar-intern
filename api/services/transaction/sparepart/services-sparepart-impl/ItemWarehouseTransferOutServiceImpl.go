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

func NewWhTransferOutImpl(transferOutRepo transactionsparepartrepository.ItemWarehouseTransferOutRepository, db *gorm.DB, redis *redis.Client) transactionsparepartservice.ItemWarehouseTransferOutService {
	return &WhTransferOutServiceImpl{
		TransferOutRepo: transferOutRepo,
		DB:              db,
		RedisClient:     redis,
	}
}

type WhTransferOutServiceImpl struct {
	TransferOutRepo transactionsparepartrepository.ItemWarehouseTransferOutRepository
	DB              *gorm.DB
	RedisClient     *redis.Client
}

// DeleteTransferOut implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) DeleteTransferOut(number int) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.DeleteTransferOut(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteTransferOutDetail implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) DeleteTransferOutDetail(number []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	var result bool
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
	for _, num := range number {
		result, err = s.TransferOutRepo.DeleteTransferOutDetail(tx, num)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// GetAllTransferOut implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) GetAllTransferOut(filter []utils.FilterCondition, dateParams map[string]string, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.GetAllTransferOut(tx, filter, dateParams, page)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAllTransferOutDetail implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) GetAllTransferOutDetail(number int, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.GetAllTransferOutDetail(tx, number, page)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetTransferOutById implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) GetTransferOutById(number int) (transactionsparepartpayloads.GetTransferOutByIdResponse, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.GetTransferOutById(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// InsertDetail implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) InsertDetail(request transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.InsertDetail(tx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// InsertDetailFromReceipt implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) InsertDetailFromReceipt(request transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.InsertDetailFromReceipt(tx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// InsertHeader implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) InsertHeader(request transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.InsertHeader(tx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SubmitTransferOut implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) SubmitTransferOut(number int) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.SubmitTransferOut(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateTransferOutDetail implements transactionsparepartservice.ItemWarehouseTransferOutService.
func (s *WhTransferOutServiceImpl) UpdateTransferOutDetail(request transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest, number int) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
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
	result, err := s.TransferOutRepo.UpdateTransferOutDetail(tx, request, number)
	if err != nil {
		return result, err
	}
	return result, nil
}
