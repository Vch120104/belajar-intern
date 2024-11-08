package masteroperationserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OperationModelMappingServiceImpl struct {
	operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository
	DB                        *gorm.DB
	RedisClient               *redis.Client // Redis client
}

func StartOperationModelMappingService(operationModelMappingRepo masteroperationrepository.OperationModelMappingRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationModelMappingService {
	return &OperationModelMappingServiceImpl{
		operationModelMappingRepo: operationModelMappingRepo,
		DB:                        db,
		RedisClient:               redisClient,
	}
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingById(id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationModelMappingById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.operationModelMappingRepo.GetOperationModelMappingLookup(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationModelMappingByBrandModelOperationCode(request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationModelMappingByBrandModelOperationCode(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMapping(req masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.SaveOperationModelMapping(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ChangeStatusOperationModelMapping(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.ChangeStatusOperationModelMapping(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingFrt(request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingFrt(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeleteOperationLevel(ids []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.operationModelMappingRepo.DeleteOperationLevel(tx, ids)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.DeactivateOperationFrt(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationFrt(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.ActivateOperationFrt(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationDocumentRequirement(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetAllOperationDocumentRequirement(tx, id, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationFrt(id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.operationModelMappingRepo.GetAllOperationFrt(tx, id, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationDocumentRequirementById(id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationDocumentRequirementById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingDocumentRequirementRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationFrtById(id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationFrtById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteroperationpayloads.OperationModelMappingFrtRequest{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationModelMappingDocumentRequirement(request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.SaveOperationModelMappingDocumentRequirement(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.DeactivateOperationDocumentRequirement(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationDocumentRequirement(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.ActivateOperationDocumentRequirement(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) SaveOperationLevel(request masteroperationpayloads.OperationLevelRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.SaveOperationLevel(tx, request)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetAllOperationLevel(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetAllOperationLevel(tx, id, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) GetOperationLevelById(id int) (masteroperationpayloads.OperationLevelByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.GetOperationLevelById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteroperationpayloads.OperationLevelByIdResponse{}, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) DeactivateOperationLevel(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.DeactivateOperationLevel(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *OperationModelMappingServiceImpl) ActivateOperationLevel(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationModelMappingRepo.ActivateOperationLevel(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
