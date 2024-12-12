package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
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

	results, repoErr := s.ItemLocationTransferRepository.GetAllItemLocationTransfer(tx, filterCondition, pages)
	if repoErr != nil {
		return results, repoErr
	}

	return results, nil
}
