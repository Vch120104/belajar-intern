package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemServiceImpl struct {
	itemRepo masteritemrepository.ItemRepository
	DB       *gorm.DB
}

func StartItemService(itemRepo masteritemrepository.ItemRepository, db *gorm.DB) masteritemservice.ItemService {
	return &ItemServiceImpl{
		itemRepo: itemRepo,
		DB:       db,
	}
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition) ([]masteritempayloads.ItemLookup, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItem(tx, filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemLookup(queryParams map[string]string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItemLookup(tx, queryParams)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetItemCode(tx, code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemId != 0 {
		_, err := s.itemRepo.GetItemById(tx, req.ItemId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.itemRepo.SaveItem(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.itemRepo.GetItemById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
