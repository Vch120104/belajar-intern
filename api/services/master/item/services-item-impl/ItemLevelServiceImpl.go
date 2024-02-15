package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"

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

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.ItemLevelId != 0 {
		_, err := s.structItemLevelRepo.GetById(tx, request.ItemLevelId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	save, err := s.structItemLevelRepo.Save(tx, request)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return save
}

func (s *ItemLevelServiceImpl) GetById(itemLevelId int) masteritemlevelpayloads.GetItemLevelResponseById {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return get
}

func (s *ItemLevelServiceImpl) GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetAll(tx, request, pages)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return get
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevelId int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	
	_, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	change_status, err := s.structItemLevelRepo.ChangeStatus(tx, itemLevelId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return change_status
}
