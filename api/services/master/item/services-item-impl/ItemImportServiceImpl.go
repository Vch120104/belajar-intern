package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemImportServiceImpl struct {
	itemImportRepo masteritemrepository.ItemImportRepository
	DB             *gorm.DB
}

// GetItemImportbyId implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GetItemImportbyId(Id int) (any, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemImportRepo.GetItemImportbyId(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) GetAllItemImport(internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.itemImportRepo.GetAllItemImport(tx, internalFilter, externalFilter, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// SaveItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) SaveItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemImportRepo.SaveItemImport(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateItemImport implements masteritemservice.ItemImportService.
func (s *ItemImportServiceImpl) UpdateItemImport(req masteritementities.ItemImport) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemImportRepo.UpdateItemImport(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func StartItemImportService(ItemImportrepo masteritemrepository.ItemImportRepository, db *gorm.DB) masteritemservice.ItemImportService {
	return &ItemImportServiceImpl{
		itemImportRepo: ItemImportrepo,
		DB:             db,
	}
}
