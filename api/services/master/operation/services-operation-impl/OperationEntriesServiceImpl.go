package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"

	"gorm.io/gorm"
)

type OperationEntriesServiceImpl struct {
	operationEntriesRepo masteroperationrepository.OperationEntriesRepository
}

func StartOperationEntriesService(operationEntriesRepo masteroperationrepository.OperationEntriesRepository) masteroperationservice.OperationEntriesService {
	return &OperationEntriesServiceImpl{
		operationEntriesRepo: operationEntriesRepo,
	}
}

func (s *OperationEntriesServiceImpl) WithTrx(trxHandle *gorm.DB) masteroperationservice.OperationEntriesService {
	s.operationEntriesRepo = s.operationEntriesRepo.WithTrx(trxHandle)
	return s
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesById(id int32) (masteroperationpayloads.OperationEntriesResponse, error) {
	results, err := s.operationEntriesRepo.GetOperationEntriesById(id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) GetOperationEntriesName(request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error) {
	results, err := s.operationEntriesRepo.GetOperationEntriesName(request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) SaveOperationEntries(req masteroperationpayloads.OperationEntriesResponse) (bool, error) {
	results, err := s.operationEntriesRepo.SaveOperationEntries(req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationEntriesServiceImpl) ChangeStatusOperationEntries(Id int) (bool, error) {
	results, err := s.operationEntriesRepo.ChangeStatusOperationEntries(Id)
	if err != nil {
		return results, err
	}
	return results, nil
}
