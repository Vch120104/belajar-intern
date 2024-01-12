package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationModelMappingServiceImpl struct {
	operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository
}

func StartOperationMappingService(operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository) masteroperationservice.OperationModelMappingService {
	return &OperationModelMappingServiceImpl{
		operationModelMappingRepo: operationModelMappingRepo,
	}
}

func (s *OperationModelMappingServiceImpl) WithTrx(trxHandle *gorm.DB) masteroperationservice.OperationModelMappingService {
	s.operationModelMappingRepo = s.operationModelMappingRepo.WithTrx(trxHandle)
	return s
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingById(id int) (masteroperationpayloads.OperationModelMappingResponse, error) {
	results, err := s.operationModelMappingRepo.GetOperationModelMappingById(id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingLookup(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	results, err := s.operationModelMappingRepo.GetOperationModelMappingLookup(filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, error) {
	results, err := s.operationModelMappingRepo.GetOperationModelMappingByBrandModelOperationCode(request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMapping(req masteroperationpayloads.OperationModelMappingResponse) (bool, error) {
	results, err := s.operationModelMappingRepo.SaveOperationModelMapping(req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ChangeStatusOperationModelMapping(Id int) (bool, error) {
	results, err := s.operationModelMappingRepo.ChangeStatusOperationModelMapping(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
