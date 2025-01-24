package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PrintGatePassServiceImpl struct {
	DB                      *gorm.DB
	PrintGatePassRepository transactionworkshoprepository.PrintGatePassRepository
	RedisClient             *redis.Client
}

func OpenPrintGatePassServiceImpl(db *gorm.DB, repository transactionworkshoprepository.PrintGatePassRepository, redisClient *redis.Client) *PrintGatePassServiceImpl {
	return &PrintGatePassServiceImpl{
		DB:                      db,
		PrintGatePassRepository: repository,
		RedisClient:             redisClient,
	}
}

func (s *PrintGatePassServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

	pages, repoErr := s.PrintGatePassRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return pages, repoErr
	}

	return pages, nil
}
