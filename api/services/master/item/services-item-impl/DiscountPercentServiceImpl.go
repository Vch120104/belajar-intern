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

type DiscountPercentServiceImpl struct {
	discountPercentRepo masteritemrepository.DiscountPercentRepository
	DB                  *gorm.DB
}

func StartDiscountPercentService(discountPercentRepo masteritemrepository.DiscountPercentRepository, db *gorm.DB) masteritemservice.DiscountPercentService {
	return &DiscountPercentServiceImpl{
		discountPercentRepo: discountPercentRepo,
		DB:                  db,
	}
}

func (s *DiscountPercentServiceImpl) GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.discountPercentRepo.GetAllDiscountPercent(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}

func (s *DiscountPercentServiceImpl) GetDiscountPercentById(Id int) masteritempayloads.DiscountPercentResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountPercentServiceImpl) SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountPercentRepo.SaveDiscountPercent(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountPercentServiceImpl) ChangeStatusDiscountPercent(Id int) (bool) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.discountPercentRepo.ChangeStatusDiscountPercent(tx, Id)
	if err != nil {
		return results
	}
	return true
}
