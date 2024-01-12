package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
)

type OperationCodeServiceImpl struct {
	operationCodeRepo masteroperationrepository.OperationCodeRepository
}

func StartOperationCodeService(operationCodeRepo masteroperationrepository.OperationCodeRepository) masteroperationservice.OperationCodeService {
	return &OperationCodeServiceImpl{
		operationCodeRepo: operationCodeRepo,
	}
}

func (s *OperationCodeServiceImpl) GetAllOperationCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.operationCodeRepo.GetAllOperationCode(filterCondition, pages)
	if err != nil {
		return pages, err
	}
	return results, nil
}

func (s *OperationCodeServiceImpl) GetOperationCodeById(id int32) (masteroperationpayloads.OperationCodeResponse, error) {
	results, err := s.operationCodeRepo.GetOperationCodeById(id)
	if err != nil {
		return results, err
	}
	return results, nil
}
