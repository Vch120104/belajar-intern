package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateServiceImpl struct {
	markupRepo masteritemrepository.MarkupRateRepository
	DB         *gorm.DB
}

func StartMarkupRateService(markupRepo masteritemrepository.MarkupRateRepository, db *gorm.DB) masteritemservice.MarkupRateService {
	return &MarkupRateServiceImpl{
		markupRepo: markupRepo,
		DB:         db,
	}
}

func (s *MarkupRateServiceImpl) GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.markupRepo.GetAllMarkupRate(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}

func (s *MarkupRateServiceImpl) GetMarkupRateById(id int) masteritempayloads.MarkupRateResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupRateById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *MarkupRateServiceImpl) SaveMarkupRate(req masteritempayloads.MarkupRateRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.SaveMarkupRate(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *MarkupRateServiceImpl) ChangeStatusMarkupRate(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.markupRepo.GetMarkupRateById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.markupRepo.ChangeStatusMarkupRate(tx, Id)
	if err != nil {
		return results
	}
	return true
}