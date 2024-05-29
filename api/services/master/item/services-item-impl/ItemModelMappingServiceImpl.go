package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

	"gorm.io/gorm"
)

type ItemModelMappingServiceImpl struct {
	ItemModelMappingRepo masteritemrepository.ItemModelMappingRepository
	DB                   *gorm.DB
}

// GetItemModelMappingByItemId implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) GetItemModelMappingByItemId(itemId int, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.ItemModelMappingRepo.GetItemModelMappingByItemId(tx, itemId, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// UpdateItemModelMapping implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) UpdateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemModelMappingRepo.UpdateItemModelMapping(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

// CreateItemModelMapping implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) CreateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemModelMappingRepo.CreateItemModelMapping(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func StartItemModelMappingService(ItemModelMappingRepo masteritemrepository.ItemModelMappingRepository, db *gorm.DB) masteritemservice.ItemModelMappingService {
	return &ItemModelMappingServiceImpl{
		ItemModelMappingRepo: ItemModelMappingRepo,
		DB:                   db,
	}
}
