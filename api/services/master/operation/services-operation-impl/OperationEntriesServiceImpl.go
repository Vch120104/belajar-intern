package masteroperationserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationEntriesServiceImpl struct {
	operationEntriesRepo masteroperationrepository.OperationEntriesRepository
	DB                   *gorm.DB
}

func StartOperationEntriesService(operationEntriesRepo masteroperationrepository.OperationEntriesRepository, db *gorm.DB) masteroperationservice.OperationEntriesService {
	return &OperationEntriesServiceImpl{
		operationEntriesRepo: operationEntriesRepo,
		DB:                   db,
	}
}

func (s *OperationEntriesServiceImpl) GetAllOperationEntries(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationEntriesRepo.GetAllOperationEntries(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesById(id int) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationEntriesRepo.GetOperationEntriesById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesName(request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationEntriesRepo.GetOperationEntriesName(tx, request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) SaveOperationEntries(req masteroperationpayloads.OperationEntriesResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationEntriesId != 0 {
		_, err := s.operationEntriesRepo.GetOperationEntriesById(tx, req.OperationEntriesId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationEntriesRepo.SaveOperationEntries(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) ChangeStatusOperationEntries(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationEntriesRepo.ChangeStatusOperationEntries(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
