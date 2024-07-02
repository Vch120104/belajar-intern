package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MarkupRateServiceImpl struct {
	markupRepo  masteritemrepository.MarkupRateRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartMarkupRateService(markupRepo masteritemrepository.MarkupRateRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.MarkupRateService {
	return &MarkupRateServiceImpl{
		markupRepo:  markupRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

func (s *MarkupRateServiceImpl) GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.markupRepo.GetAllMarkupRate(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *MarkupRateServiceImpl) GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.markupRepo.GetMarkupRateById(tx, id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *MarkupRateServiceImpl) SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.MarkupRateId != 0 {
		_, err := s.markupRepo.GetMarkupRateById(tx, req.MarkupRateId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.markupRepo.SaveMarkupRate(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *MarkupRateServiceImpl) ChangeStatusMarkupRate(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.markupRepo.GetMarkupRateById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.markupRepo.ChangeStatusMarkupRate(tx, Id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return true, nil
}

func (s *MarkupRateServiceImpl) GetMarkupRateByMarkupMasterAndOrderType(MarkupMasterId int, OrderTypeId int) ([]masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.markupRepo.GetMarkupRateByMarkupMasterAndOrderType(tx, MarkupMasterId, OrderTypeId)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}
