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

type ItemLocationServiceImpl struct {
	ItemLocationRepo masteritemrepository.ItemLocationRepository
	DB               *gorm.DB
}

func StartItemLocationService(ItemLocationRepo masteritemrepository.ItemLocationRepository, db *gorm.DB) masteritemservice.ItemLocationService {
	return &ItemLocationServiceImpl{
		ItemLocationRepo: ItemLocationRepo,
		DB:               db,
	}
}

func (s *ItemLocationServiceImpl) GetAllItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.ItemLocationRepo.GetAllItemLocation(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemLocationServiceImpl) SaveItemLocation(req masteritempayloads.ItemLocationRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemLocationRepo.SaveItemLocation(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}
