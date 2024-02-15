package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassServiceImpl struct {
	itemRepo masteritemrepository.ItemClassRepository
	DB       *gorm.DB
}

func StartItemClassService(itemRepo masteritemrepository.ItemClassRepository, db *gorm.DB) masteritemservice.ItemClassService {
	return &ItemClassServiceImpl{
		itemRepo: itemRepo,
		DB:       db,
	}
}

func (s *ItemClassServiceImpl) GetAllItemClass(filterCondition []utils.FilterCondition) []map[string]interface{} {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItemClass(tx, filterCondition)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemClassServiceImpl) GetItemClassById(Id int) masteritempayloads.ItemClassResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemClassById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}

func (s *ItemClassServiceImpl) SaveItemClass(req masteritempayloads.ItemClassResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemClassId != 0 {
		_, err := s.itemRepo.GetItemClassById(tx, req.ItemClassId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.itemRepo.SaveItemClass(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemClassServiceImpl) ChangeStatusItemClass(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.itemRepo.GetItemClassById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.itemRepo.ChangeStatusItemClass(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
