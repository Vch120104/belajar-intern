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

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemInquiryServiceImpl struct {
	ItemInquiryRepo transactionsparepartrepository.ItemInquiryRepository
	DB              *gorm.DB
	redis           *redis.Client
}

func StartItemInquiryService(itemInquiryRepo transactionsparepartrepository.ItemInquiryRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.ItemInquiryService {
	return &ItemInquiryServiceImpl{
		ItemInquiryRepo: itemInquiryRepo,
		DB:              db,
		redis:           rdb,
	}
}

func (i *ItemInquiryServiceImpl) GetAllItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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

	result, err := i.ItemInquiryRepo.GetAllItemInquiry(tx, filterCondition, pages)
	if err != nil {
		return pages, err
	}
	return result, nil
}

func (i *ItemInquiryServiceImpl) GetByIdItemInquiry(filter transactionsparepartpayloads.ItemInquiryGetByIdFilter) (transactionsparepartpayloads.ItemInquiryGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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

	result, err := i.ItemInquiryRepo.GetByIdItemInquiry(tx, filter)
	if err != nil {
		return result, err
	}
	return result, nil
}
