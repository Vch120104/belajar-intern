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

type QualityControlServiceImpl struct {
	QualityControlRepository transactionworkshoprepository.QualityControlRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenQualityControlServiceImpl(QualityControlRepo transactionworkshoprepository.QualityControlRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.QualityControlService {
	return &QualityControlServiceImpl{
		QualityControlRepository: QualityControlRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *QualityControlServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.QualityControlRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (s *QualityControlServiceImpl) GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.GetById(tx, id, filterCondition, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlServiceImpl) Qcpass(id int, iddet int) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.Qcpass(tx, id, iddet)
	if repoErr != nil {

		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlServiceImpl) Reorder(id int, iddet int, payload transactionworkshoppayloads.QualityControlReorder) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.Reorder(tx, id, iddet, payload)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}
