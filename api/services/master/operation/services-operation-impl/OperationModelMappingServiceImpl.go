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

type OperationModelMappingServiceImpl struct {
	operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository
	DB                        *gorm.DB
}

func StartOperationModelMappingService(operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository, db *gorm.DB) masteroperationservice.OperationModelMappingService {
	return &OperationModelMappingServiceImpl{
		operationModelMappingRepo: operationModelMappingRepo,
		DB:                        db,
	}
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingById(id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetOperationModelMappingById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.operationModelMappingRepo.GetOperationModelMappingLookup(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetOperationModelMappingByBrandModelOperationCode(tx, request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMapping(req masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.SaveOperationModelMapping(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ChangeStatusOperationModelMapping(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.ChangeStatusOperationModelMapping(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingFrt(request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingFrt(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationFrt(id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.DeactivateOperationFrt(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationFrt(id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.ActivateOperationFrt(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationDocumentRequirement(id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetAllOperationDocumentRequirement(tx, id, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationFrt(id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetAllOperationFrt(tx, id, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationDocumentRequirementById(id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetOperationDocumentRequirementById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingDocumentRequirementRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationFrtById(id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.GetOperationFrtById(tx, id)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingFrtRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingDocumentRequirement(request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingDocumentRequirement(tx, request)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationDocumentRequirement(id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.DeactivateOperationDocumentRequirement(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationDocumentRequirement(id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationModelMappingRepo.ActivateOperationDocumentRequirement(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}
