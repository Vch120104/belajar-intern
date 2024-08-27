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

type QualityControlBodyshopServiceImpl struct {
	QualityControlRepository transactionbodyshoprepository.QualityControlBodyshopRepository
	DB                       *gorm.DB
	RedisClient              *redis.Client // Redis client
}

func OpenQualityControlBodyshopServiceImpl(QualityControlRepo transactionbodyshoprepository.QualityControlBodyshopRepository, db *gorm.DB, redisClient *redis.Client) transactionbodyshopservice.QualityControlBodyshopService {
	return &QualityControlBodyshopServiceImpl{
		QualityControlRepository: QualityControlRepo,
		DB:                       db,
		RedisClient:              redisClient,
	}
}

func (s *QualityControlBodyshopServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	results, totalPages, totalRows, repoErr := s.QualityControlRepository.GetAll(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(results, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (s *QualityControlBodyshopServiceImpl) GetById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.GetById(tx, id, filterCondition, pages)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlBodyshopServiceImpl) Qcpass(id int, iddet int) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.Qcpass(tx, id, iddet)
	if repoErr != nil {

		return result, repoErr
	}

	return result, nil
}

func (s *QualityControlBodyshopServiceImpl) Reorder(id int, iddet int, payload transactionbodyshoppayloads.QualityControlReorder) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollbackTrx(tx)

	result, repoErr := s.QualityControlRepository.Reorder(tx, id, iddet, payload)
	if repoErr != nil {
		return result, repoErr
	}

	return result, nil
}
