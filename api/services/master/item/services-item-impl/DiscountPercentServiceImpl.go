package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
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

// func (s *DiscountPercentServiceImpl) GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error) {
// 	result, err := s.discountPercentRepo.GetDiscountPercentById(Id)
// 	if err != nil {
// 		return result, err
// 	}
// 	return result, nil
// }

// func (s *DiscountPercentServiceImpl) SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, error) {
// 	results, err := s.discountPercentRepo.SaveDiscountPercent(req)
// 	if err != nil {
// 		return results, err
// 	}
// 	return results, nil
// }

// func (s *DiscountPercentServiceImpl) ChangeStatusDiscountPercent(Id int) (bool, error) {
// 	results, err := s.discountPercentRepo.ChangeStatusDiscountPercent(Id)
// 	if err != nil {
// 		return results, err
// 	}
// 	return results, nil
// }
