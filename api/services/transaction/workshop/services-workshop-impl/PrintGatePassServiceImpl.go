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
	PrintGatePassRepository transactionworkshoprepository.PrintGatePassRepository
	DB                      *gorm.DB
	RedisClient             *redis.Client
}

func OpenPrintGatePassServiceImpl(PrintGatePassRepo transactionworkshoprepository.PrintGatePassRepository, db *gorm.DB, redisClient *redis.Client) *PrintGatePassServiceImpl {
	return &PrintGatePassServiceImpl{
		PrintGatePassRepository: PrintGatePassRepo,
		DB:                      db,
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

func (s *PrintGatePassServiceImpl) PrintById(id int) ([]byte, *exceptions.BaseErrorResponse) {
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

	// Ambil data dari repository
	dataMap, repoErr := s.PrintGatePassRepository.PrintById(tx, id)
	if repoErr != nil {
		return nil, repoErr
	}

	// Generate PDF menggunakan `GeneratePDFGatePass`
	pdfBytes, pdfErr := utils.GeneratePDFGatePass(dataMap)
	if pdfErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate PDF",
			Err:        pdfErr,
		}
	}

	return pdfBytes, nil
}
