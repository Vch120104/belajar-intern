package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupServiceImpl struct {
	operationGroupRepo masteroperationrepository.OperationGroupRepository
}

func StartOperationGroupService(operationGroupRepo masteroperationrepository.OperationGroupRepository) masteroperationservice.OperationGroupService {
	return &OperationGroupServiceImpl{
		operationGroupRepo: operationGroupRepo,
	}
}
func (s *OperationGroupServiceImpl) WithTrx(trxHandle *gorm.DB) masteroperationservice.OperationGroupService {
	s.operationGroupRepo = s.operationGroupRepo.WithTrx(trxHandle)
	return s
}

func (s *OperationGroupServiceImpl) GetAllOperationGroupIsActive() ([]masteroperationpayloads.OperationGroupResponse, error) {
	results, err := s.operationGroupRepo.GetAllOperationGroupIsActive()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupById(id int) (masteroperationpayloads.OperationGroupResponse, error) {
	results, err := s.operationGroupRepo.GetOperationGroupById(id)

	if err != nil {
		return masteroperationpayloads.OperationGroupResponse{}, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupByCode(Code string) (masteroperationpayloads.OperationGroupResponse, error) {
	results, err := s.operationGroupRepo.GetOperationGroupByCode(Code)
	if err != nil {
		return masteroperationpayloads.OperationGroupResponse{}, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.operationGroupRepo.GetAllOperationGroup(filterCondition, pages)
	if err != nil {
		return pages, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) ChangeStatusOperationGroup(oprId int) (bool, error) {
	results, err := s.operationGroupRepo.ChangeStatusOperationGroup(oprId)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) SaveOperationGroup(req masteroperationpayloads.OperationGroupResponse) (bool, error) {
	results, err := s.operationGroupRepo.SaveOperationGroup(req)
	if err != nil {
		return false, err
	}
	return results, nil
}
