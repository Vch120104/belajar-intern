package masteroperationserviceimpl

import (
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	// "after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationSectionServiceImpl struct {
	operationSectionRepo masteroperationrepository.OperationSectionRepository
}

func StartOperationSectionService(operationSectionRepo masteroperationrepository.OperationSectionRepository) masteroperationservice.OperationSectionService {
	return &OperationSectionServiceImpl{
		operationSectionRepo: operationSectionRepo,
	}
}

func (r *OperationSectionServiceImpl) WithTrx(trxHandle *gorm.DB) masteroperationservice.OperationSectionService {
	r.operationSectionRepo = r.operationSectionRepo.WithTrx(trxHandle)
	return r
}

func (s *OperationSectionServiceImpl) GetAllOperationSectionList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.operationSectionRepo.GetAllOperationSectionList(filterCondition, pages)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *OperationSectionServiceImpl) GetSectionCodeByGroupId(GroupId string) ([]masteroperationpayloads.OperationSectionCodeResponse, error) {
	results, err := s.operationSectionRepo.GetSectionCodeByGroupId(GroupId)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionName(group_id int, section_code string) (masteroperationpayloads.OperationSectionNameResponse, error) {
	results, err := s.operationSectionRepo.GetOperationSectionName(group_id, section_code)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *OperationSectionServiceImpl) SaveOperationSection(req masteroperationpayloads.OperationSectionRequest) (bool, error) {
	results, err := s.operationSectionRepo.SaveOperationSection(req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionById(id int) (masteroperationpayloads.OperationSectionResponse, error) {
	results, err := s.operationSectionRepo.GetOperationSectionById(id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetAllOperationSection() ([]masteroperationpayloads.OperationSectionResponse, error) {
	results, err := s.operationSectionRepo.GetAllOperationSection()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) ChangeStatusOperationSection(Id int) (bool, error) {
	results, err := s.operationSectionRepo.ChangeStatusOperationSection(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
