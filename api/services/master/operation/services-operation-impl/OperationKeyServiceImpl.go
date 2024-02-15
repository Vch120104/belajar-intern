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

type OperationKeyServiceImpl struct {
	operationKeyRepo masteroperationrepository.OperationKeyRepository
	DB               *gorm.DB
}

func StartOperationKeyService(operationKeyRepo masteroperationrepository.OperationKeyRepository, db *gorm.DB) masteroperationservice.OperationKeyService {
	return &OperationKeyServiceImpl{
		operationKeyRepo: operationKeyRepo,
		DB:               db,
	}
}

func (s *OperationKeyServiceImpl) GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetAllOperationKeyList(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationKeyServiceImpl) GetOperationKeyById(id int) masteroperationpayloads.OperationkeyListResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetOperationKeyById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationKeyServiceImpl) GetOperationKeyName(req masteroperationpayloads.OperationKeyRequest) masteroperationpayloads.OperationKeyNameResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetOperationKeyName(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))

	}
	return results
}

func (s *OperationKeyServiceImpl) SaveOperationKey(req masteroperationpayloads.OperationKeyResponse) bool {

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationKeyId != 0 {
		_, err := s.operationKeyRepo.GetOperationKeyById(tx, req.OperationKeyId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.operationKeyRepo.SaveOperationKey(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationKeyServiceImpl) ChangeStatusOperationKey(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.ChangeStatusOperationKey(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
