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

type MarkupMasterServiceImpl struct {
	markupRepo masteritemrepository.MarkupMasterRepository
	DB         *gorm.DB
}

func StartMarkupMasterService(markupRepo masteritemrepository.MarkupMasterRepository, db *gorm.DB) masteritemservice.MarkupMasterService {
	return &MarkupMasterServiceImpl{
		markupRepo: markupRepo,
		DB:         db,
	}
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupMasterList(tx, filter, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterById(id int) masteritempayloads.MarkupMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupMasterById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *MarkupMasterServiceImpl) SaveMarkupMaster(req masteritempayloads.MarkupMasterResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.MarkupMasterId != 0 {
		_, err := s.markupRepo.GetMarkupMasterById(tx, req.MarkupMasterId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.markupRepo.SaveMarkupMaster(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
func (s *MarkupMasterServiceImpl) ChangeStatusMasterMarkupMaster(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.markupRepo.GetMarkupMasterById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.markupRepo.ChangeStatusMasterMarkupMaster(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
func (s *MarkupMasterServiceImpl) GetMarkupMasterByCode(markupCode string) masteritempayloads.MarkupMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.markupRepo.GetMarkupMasterByCode(tx, markupCode)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result

}
