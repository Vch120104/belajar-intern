package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemLocationTransferServiceImpl struct {
	ItemLocationTransferRepository transactionsparepartrepository.ItemLocationTransferRepository
	DB                             *gorm.DB
}

func NewItemLocationTransferServiceImpl(
	itemLocationTransferRepo transactionsparepartrepository.ItemLocationTransferRepository,
	db *gorm.DB,
) transactionsparepartservice.ItemLocationTransferService {
	return &ItemLocationTransferServiceImpl{
		ItemLocationTransferRepository: itemLocationTransferRepo,
		DB:                             db,
	}
}

func (s *ItemLocationTransferServiceImpl) GetAllItemLocationTransfer(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.GetAllItemLocationTransfer(tx, filterCondition, pages)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) GetItemLocationTransferById(id int) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.GetItemLocationTransferById(tx, id)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) InsertItemLocationTransfer(request transactionsparepartpayloads.InsertItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.InsertItemLocationTransfer(tx, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) UpdateItemLocationTransfer(id int, request transactionsparepartpayloads.UpdateItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.UpdateItemLocationTransfer(tx, id, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) AcceptItemLocationTransfer(id int, request transactionsparepartpayloads.AcceptItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.AcceptItemLocationTransfer(tx, id, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) RejectItemLocationTransfer(id int, request transactionsparepartpayloads.RejectItemLocationTransferRequest) (transactionsparepartpayloads.GetItemLocationTransferByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.RejectItemLocationTransfer(tx, id, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) DeleteItemLocationTransfer(id int) (bool, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.DeleteItemLocationTransfer(tx, id)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) InsertItemLocationTransferDetail(request transactionsparepartpayloads.InsertItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.InsertItemLocationTransferDetail(tx, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) UpdateItemLocationTransferDetail(id int, request transactionsparepartpayloads.UpdateItemLocationTransferDetailRequest) (transactionsparepartpayloads.GetItemLocationTransferDetailByIdResponse, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.UpdateItemLocationTransferDetail(tx, id, request)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

func (s *ItemLocationTransferServiceImpl) DeleteItemLocationTransferDetail(ids []int) (bool, *exceptions.BaseErrorResponse) {
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

	response, responseErr := s.ItemLocationTransferRepository.DeleteItemLocationTransferDetail(tx, ids)
	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}
