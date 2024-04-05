package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
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

func (s *IncentiveMasterServiceImpl) GetAllIncentiveMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.IncentiveMasterRepo.GetAllIncentiveMaster(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *IncentiveMasterServiceImpl) GetIncentiveMasterById(id int) (masterpayloads.IncentiveMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveMasterRepo.GetIncentiveMasterById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *IncentiveMasterServiceImpl) SaveIncentiveMaster(req masterpayloads.IncentiveMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.IncentiveMasterRepo.SaveIncentiveMaster(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *IncentiveMasterServiceImpl) ChangeStatusIncentiveMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.IncentiveMasterRepo.GetIncentiveMasterById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.IncentiveMasterRepo.ChangeStatusIncentiveMaster(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
