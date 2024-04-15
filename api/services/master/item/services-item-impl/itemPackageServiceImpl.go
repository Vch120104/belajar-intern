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

type ItemPackageServiceImpl struct {
	ItemPackageRepo masteritemrepository.ItemPackageRepository
	DB              *gorm.DB
}

func StartItemPackageService(ItemPackageRepo masteritemrepository.ItemPackageRepository, db *gorm.DB) masteritemservice.ItemPackageService {
	return &ItemPackageServiceImpl{
		ItemPackageRepo: ItemPackageRepo,
		DB:              db,
	}
}

func (s *ItemPackageServiceImpl) GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.ItemPackageRepo.GetAllItemPackage(tx, internalFilterCondition, externalFilterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemPackageServiceImpl) GetItemPackageById(Id int) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageRepo.GetItemPackageById(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) SaveItemPackage(req masteritempayloads.SaveItemPackageRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemPackageId != 0 {
		_, err := s.ItemPackageRepo.GetItemPackageById(tx, req.ItemPackageId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.ItemPackageRepo.SaveItemPackage(tx, req)

	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) ChangeStatusItemPackage(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ItemPackageRepo.GetItemPackageById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.ItemPackageRepo.ChangeStatusItemPackage(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
