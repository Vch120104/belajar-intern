package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
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

func (s *DiscountPercentServiceImpl) GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.discountPercentRepo.GetAllDiscountPercent(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *DiscountPercentServiceImpl) GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountPercentServiceImpl) SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	if req.DiscountPercentId != 0 {
		_, err := s.discountPercentRepo.GetDiscountPercentById(tx, req.DiscountPercentId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.discountPercentRepo.SaveDiscountPercent(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *DiscountPercentServiceImpl) ChangeStatusDiscountPercent(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.discountPercentRepo.GetDiscountPercentById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.discountPercentRepo.ChangeStatusDiscountPercent(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
