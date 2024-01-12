package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationKeyServiceImpl struct {
	operationKeyRepo masteroperationrepository.OperationKeyRepository
}

func StartOperationKeyService(operationKeyRepo masteroperationrepository.OperationKeyRepository) masteroperationservice.OperationKeyService {
	return &OperationKeyServiceImpl{
		operationKeyRepo: operationKeyRepo,
	}
}

func (s *OperationKeyServiceImpl) WithTrx(trxHandle *gorm.DB) masteroperationservice.OperationKeyService {
	s.operationKeyRepo = s.operationKeyRepo.WithTrx(trxHandle)
	return s
}

func (s *OperationKeyServiceImpl) GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.operationKeyRepo.GetAllOperationKeyList(filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyById(id int) (masteroperationpayloads.OperationKeyResponse, error) {
	results, err := s.operationKeyRepo.GetOperationKeyById(id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) GetOperationKeyName(req masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error) {
	results, err := s.operationKeyRepo.GetOperationKeyName(req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) SaveOperationKey(req masteroperationpayloads.OperationKeyResponse) (bool, error) {
	results, err := s.operationKeyRepo.SaveOperationKey(req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationKeyServiceImpl) ChangeStatusOperationKey(Id int) (bool, error) {
	results, err := s.operationKeyRepo.ChangeStatusOperationKey(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
