package transactionworkshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceWorkshopServiceImpl struct {
	ServiceWorkshopRepository transactionworkshoprepository.ServiceWorkshopRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func OpenServiceWorkshopServiceImpl(ServiceWorkshopRepo transactionworkshoprepository.ServiceWorkshopRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.ServiceWorkshopService {
	return &ServiceWorkshopServiceImpl{
		ServiceWorkshopRepository: ServiceWorkshopRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *ServiceWorkshopServiceImpl) GetAllByTechnicianWO(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ServiceWorkshopDetailResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.ServiceWorkshopRepository.GetAllByTechnicianWO(tx, idTech, idSysWo, filterCondition, pages)
	if repoErr != nil {
		return transactionworkshoppayloads.ServiceWorkshopDetailResponse{}, repoErr
	}

	return results, nil
}

func (s *ServiceWorkshopServiceImpl) StartService(idAlloc int, idSysWo int, idServLog int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Start the service
	start, err := s.ServiceWorkshopRepository.StartService(tx, idAlloc, idSysWo, idServLog, companyId)
	if err != nil {
		return false, err
	}

	return start, nil
}

func (s *ServiceWorkshopServiceImpl) PendingService(idAlloc int, idSysWo int, idServLog int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Pending the service
	pending, err := s.ServiceWorkshopRepository.PendingService(tx, idAlloc, idSysWo, idServLog, companyId)
	if err != nil {
		return false, err
	}

	return pending, nil
}
