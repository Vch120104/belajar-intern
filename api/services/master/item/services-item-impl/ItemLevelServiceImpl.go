package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLevelServiceImpl struct {
	structItemLevelRepo masteritemlevelrepo.ItemLevelRepository
	DB                  *gorm.DB
}

func StartItemLevelService(itemlevelrepo masteritemlevelrepo.ItemLevelRepository, db *gorm.DB) masteritemlevelservice.ItemLevelService {
	return &ItemLevelServiceImpl{
		structItemLevelRepo: itemlevelrepo,
		DB:                  db,
	}
}

// GetItemLevelLookUp implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetItemLevelLookUp(tx, filter, pages, itemClassId)

	if err != nil {
		return get, err
	}

	return get, nil
}

// GetItemLevelDropDown implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetItemLevelDropDown(tx, itemLevel)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.ItemLevelId != 0 {
		_, err := s.structItemLevelRepo.GetById(tx, request.ItemLevelId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.structItemLevelRepo.Save(tx, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *ItemLevelServiceImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetAll(tx, filter, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevelId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		return false, err
	}

	change_status, err := s.structItemLevelRepo.ChangeStatus(tx, itemLevelId)

	if err != nil {
		return change_status, err
	}

	return true, nil
}
