package transactionworkshopserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
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
func (s *LicenseOwnerChangeServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)
	results, totalPages, totalRows, repoErr := s.LicenseOwncerChangeRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)
	return paginatedData, totalPages, totalRows, nil
}
