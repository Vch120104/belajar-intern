package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountServiceImpl struct {
	discountRepo masterrepository.DiscountRepository
	DB           *gorm.DB
}

func StartDiscountService(discountRepo masterrepository.DiscountRepository, db *gorm.DB) masterservice.DiscountService {
	return &DiscountServiceImpl{
		discountRepo: discountRepo,
		DB:           db,
	}
}

func (s *DiscountServiceImpl) GetAllDiscountIsActive() []masterpayloads.DiscountResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetAllDiscountIsActive(tx)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountServiceImpl) GetDiscountById(id int) masterpayloads.DiscountResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetDiscountById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountServiceImpl) GetDiscountByCode(Code string) masterpayloads.DiscountResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetDiscountByCode(tx, Code)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountServiceImpl) GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.GetAllDiscount(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountServiceImpl) ChangeStatusDiscount(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountRepo.ChangeStatusDiscount(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *DiscountServiceImpl) SaveDiscount(req masterpayloads.DiscountResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.DiscountCodeId != 0 {
		_, err := s.discountRepo.GetDiscountById(tx, req.DiscountCodeId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}
	
	results, err := s.discountRepo.SaveDiscount(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
