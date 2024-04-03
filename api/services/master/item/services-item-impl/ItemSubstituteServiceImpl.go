package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemSubstituteServiceImpl struct {
	itemSubstituteRepo masteritemrepository.ItemSubstituteRepository
	Db                 *gorm.DB
}

func StartItemSubstituteService(itemSubstituteRepo masteritemrepository.ItemSubstituteRepository, db *gorm.DB) masteritemservice.ItemSubstituteService {
	return &ItemSubstituteServiceImpl{
		itemSubstituteRepo: itemSubstituteRepo,
		Db:                 db,
	}
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemSubstituteRepo.GetAllItemSubstitute(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return results
}

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstitute(id int) masteritempayloads.ItemSubstitutePayloads {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetByIdItemSubstitute(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstituteDetail(pages pagination.Pagination, id int) pagination.Pagination {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetAllItemSubstituteDetail(tx, pages, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstituteDetail(id int) masteritempayloads.ItemSubstituteDetailGetPayloads {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetByIdItemSubstituteDetail(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) bool {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.SaveItemSubstitute(tx, req)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) bool {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.SaveItemSubstituteDetail(tx, req,id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) ChangeStatusItemSubstitute(id int) bool {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.ChangeStatusItemSubstitute(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) DeactivateItemSubstituteDetail(id string) bool {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.DeactivateItemSubstituteDetail(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *ItemSubstituteServiceImpl) ActivateItemSubstituteDetail(id string) bool {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.ActivateItemSubstituteDetail(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}
