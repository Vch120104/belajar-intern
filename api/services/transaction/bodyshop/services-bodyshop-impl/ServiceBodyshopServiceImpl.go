package transactionbodyshopserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	transactionbodyshoprepository "after-sales/api/repositories/transaction/bodyshop"
	transactionbodyshopservice "after-sales/api/services/transaction/bodyshop"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceBodyshopServiceImpl struct {
	ServiceBodyshopRepository transactionbodyshoprepository.ServiceBodyshopRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func OpenServiceBodyshopServiceImpl(ServiceBodyshopRepo transactionbodyshoprepository.ServiceBodyshopRepository, db *gorm.DB, redisClient *redis.Client) transactionbodyshopservice.ServiceBodyshopService {
	return &ServiceBodyshopServiceImpl{
		ServiceBodyshopRepository: ServiceBodyshopRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *ServiceBodyshopServiceImpl) GetAllByTechnicianWOBodyshop(idTech int, idSysWo int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.ServiceBodyshopDetailResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, repoErr := s.ServiceBodyshopRepository.GetAllByTechnicianWOBodyshop(tx, idTech, idSysWo, filterCondition, pages)
	if repoErr != nil {
		return transactionbodyshoppayloads.ServiceBodyshopDetailResponse{}, repoErr
	}

	return results, nil
}

func (s *ServiceBodyshopServiceImpl) StartService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Start the service
	start, err := s.ServiceBodyshopRepository.StartService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return start, nil
}

func (s *ServiceBodyshopServiceImpl) PendingService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Pending the service
	pending, err := s.ServiceBodyshopRepository.PendingService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return pending, nil
}

func (s *ServiceBodyshopServiceImpl) TransferService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Transfer the service
	transfer, err := s.ServiceBodyshopRepository.TransferService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return transfer, nil
}

func (s *ServiceBodyshopServiceImpl) StopService(idAlloc int, idSysWo int, companyId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	// Stop the service
	stop, err := s.ServiceBodyshopRepository.StopService(tx, idAlloc, idSysWo, companyId)
	if err != nil {
		return false, err
	}

	return stop, nil
}
