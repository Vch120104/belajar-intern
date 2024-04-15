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

func (s *OperationKeyServiceImpl) GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetAllOperationKeyList(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyById(id int) (masteroperationpayloads.OperationkeyListResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetOperationKeyById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyName(req masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.GetOperationKeyName(tx, req)
	if err != nil {
		return results, err

	}
	return results, nil
}

func (s *OperationKeyServiceImpl) SaveOperationKey(req masteroperationpayloads.OperationKeyResponse) (bool, *exceptionsss_test.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationKeyId != 0 {
		_, err := s.operationKeyRepo.GetOperationKeyById(tx, req.OperationKeyId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationKeyRepo.SaveOperationKey(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) ChangeStatusOperationKey(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationKeyRepo.ChangeStatusOperationKey(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
