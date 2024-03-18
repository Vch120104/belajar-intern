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

type OperationCodeServiceImpl struct {
	operationCodeRepo masteroperationrepository.OperationCodeRepository
	DB                *gorm.DB
}

func StartOperationCodeService(operationCodeRepo masteroperationrepository.OperationCodeRepository, db *gorm.DB) masteroperationservice.OperationCodeService {
	return &OperationCodeServiceImpl{
		operationCodeRepo: operationCodeRepo,
		DB:                db,
	}
}

func (s *OperationCodeServiceImpl) GetAllOperationCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationCodeRepo.GetAllOperationCode(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationCodeServiceImpl) GetOperationCodeById(id int) masteroperationpayloads.OperationCodeResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationCodeRepo.GetOperationCodeById(tx, id)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return results
}

func (s *OperationCodeServiceImpl) SaveOperationCode(req masteroperationpayloads.OperationCodeSave) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.operationCodeRepo.SaveOperationCode(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}

func (s *OperationCodeServiceImpl) ChangeStatusOperationCode(id int) bool {
	tx := s.DB.Statement.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.operationCodeRepo.ChangeStatusItemSubstitute(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}
