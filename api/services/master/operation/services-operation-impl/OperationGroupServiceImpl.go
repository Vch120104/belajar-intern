package masteroperationserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupServiceImpl struct {
	operationGroupRepo masteroperationrepository.OperationGroupRepository
	DB                 *gorm.DB
}

func StartOperationGroupService(operationGroupRepo masteroperationrepository.OperationGroupRepository, db *gorm.DB) masteroperationservice.OperationGroupService {
	return &OperationGroupServiceImpl{
		operationGroupRepo: operationGroupRepo,
		DB:                 db,
	}
}

func (s *OperationGroupServiceImpl) GetAllOperationGroupIsActive() []masteroperationpayloads.OperationGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.operationGroupRepo.GetAllOperationGroupIsActive(tx)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}

	return get
}

func (s *OperationGroupServiceImpl) GetOperationGroupById(id int) masteroperationpayloads.OperationGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationGroupRepo.GetOperationGroupById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationGroupServiceImpl) GetOperationGroupByCode(Code string) masteroperationpayloads.OperationGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationGroupRepo.GetOperationGroupByCode(tx, Code)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationGroupServiceImpl) GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationGroupRepo.GetAllOperationGroup(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationGroupServiceImpl) ChangeStatusOperationGroup(oprId int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.operationGroupRepo.GetOperationGroupById(tx, oprId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.operationGroupRepo.ChangeStatusOperationGroup(tx, oprId)
	if err != nil {
		return results
	}
	return true
}

func (s *OperationGroupServiceImpl) SaveOperationGroup(req masteroperationpayloads.OperationGroupResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationGroupId != 0 {
		_, err := s.operationGroupRepo.GetOperationGroupById(tx, req.OperationGroupId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.operationGroupRepo.SaveOperationGroup(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
