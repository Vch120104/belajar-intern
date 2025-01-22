package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LicenseOwnerChangeServiceImpl struct {
	LicenseOwncerChangeRepository transactionworkshoprepository.LicenseOwncerChangeRepository
	DB                            *gorm.DB
	RedisClient                   *redis.Client
}

func OpenLicenseOwnerChangeServiceImpl(LicenseOwnerChangeRepo transactionworkshoprepository.LicenseOwncerChangeRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.LicenseOwnerChangeService {
	return &LicenseOwnerChangeServiceImpl{
		LicenseOwncerChangeRepository: LicenseOwnerChangeRepo,
		DB:                            db,
		RedisClient:                   redisClient,
	}
}

// GetAll implements transactionworkshopservice.LicenseOwnerChangeService.
func (s *LicenseOwnerChangeServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

    // Calling the repository method that returns paginated data
    pages, repoErr := s.LicenseOwncerChangeRepository.GetAll(tx, filterCondition, pages)
    if repoErr != nil {
        return pages, repoErr
    }

    return pages, nil
}

// GetHistoryByChassisNumber implements transactionworkshopservice.LicenseOwnerChangeService.
func (s *LicenseOwnerChangeServiceImpl) GetHistoryByChassisNumber(chassisNumber string, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

    // Calling the repository method that returns paginated data
    pages, repoErr := s.LicenseOwncerChangeRepository.GetHistoryByChassisNumber(chassisNumber, tx, filterCondition, pages)
    if repoErr != nil {
        return pages, repoErr
    }

    return pages, nil
}
