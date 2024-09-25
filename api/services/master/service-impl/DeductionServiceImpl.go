package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"

	// "after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DeductionServiceImpl struct {
	deductionrepo masterrepository.DeductionRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartDeductionService(deductionRepo masterrepository.DeductionRepository, db *gorm.DB, redisClient *redis.Client) masterservice.DeductionService {
	return &DeductionServiceImpl{
		deductionrepo: deductionRepo,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *DeductionServiceImpl) GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetAllDeduction(tx, filterCondition, pages)
	if dbErr != nil {
		// Handle error from the database operation
		return pagination.Pagination{}, dbErr
	}

	defer helper.CommitOrRollback(tx, dbErr)
	return result, nil
}

func (s *DeductionServiceImpl) GetByIdDeductionDetail(Id int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse) {

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetByIdDeductionDetail(tx, Id)
	if dbErr != nil {
		// Handle error
		return result, dbErr // Return the existing BaseErrorResponse
	}

	defer helper.CommitOrRollback(tx, dbErr)
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionList(req masterpayloads.DeductionListResponse) (masterentities.DeductionList, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.deductionrepo.SaveDeductionList(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionDetail(req masterpayloads.DeductionDetailResponse, id int) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.deductionrepo.SaveDeductionDetail(tx, req, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetDeductionById(Id int, paginate pagination.Pagination) (masterpayloads.DeductionListResponse, *exceptions.BaseErrorResponse) {

	// If data is not available in cache, fetch it from the database
	tx := s.DB.Begin()
	result, dbErr := s.deductionrepo.GetDeductionById(tx, Id, paginate)
	defer helper.CommitOrRollback(tx, dbErr)
	if dbErr != nil {
		// Handle error
		return masterpayloads.DeductionListResponse{}, dbErr
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetAllDeductionDetail(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	detail_result, detail_err := s.deductionrepo.GetAllDeductionDetail(tx, pages, Id)
	defer helper.CommitOrRollback(tx, detail_err)

	if detail_err != nil {
		return detail_result, detail_err
	}
	return detail_result, nil
}

func (s *DeductionServiceImpl) ChangeStatusDeduction(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.deductionrepo.ChangeStatusDeduction(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DeductionServiceImpl) UpdateDeductionDetail(id int, req masterpayloads.DeductionDetailUpdate) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.deductionrepo.UpdateDeductionDetail(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.DeductionDetail{}, err
	}
	return result, nil
}
