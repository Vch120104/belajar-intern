package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveMasterServiceImpl struct {
	IncentiveMasterRepo masterrepository.IncentiveMasterRepository
	DB                  *gorm.DB
}

func StartIncentiveMasterService(IncentiveMasterRepo masterrepository.IncentiveMasterRepository, db *gorm.DB) masterservice.IncentiveMasterService {
	return &IncentiveMasterServiceImpl{
		IncentiveMasterRepo: IncentiveMasterRepo,
		DB:                  db,
	}
}

func (s *IncentiveMasterServiceImpl) GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.IncentiveMasterRepo.GetAllIncentiveMaster(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}

func (s *IncentiveMasterServiceImpl) GetIncentiveMasterById(id int) masterpayloads.IncentiveMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveMasterRepo.GetIncentiveMasterById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveMasterServiceImpl) SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveMasterRepo.SaveIncentiveMaster(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *IncentiveMasterServiceImpl) ChangeStatusIncentiveMaster(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.IncentiveMasterRepo.GetIncentiveMasterById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.IncentiveMasterRepo.ChangeStatusIncentiveMaster(tx, Id)
	if err != nil {
		return results
	}
	return true
}
