package masteroperationserviceimpl

import (
	// "after-sales/api/exceptions"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"net/http"

	// "after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OperationSectionServiceImpl struct {
	operationSectionRepo masteroperationrepository.OperationSectionRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func StartOperationSectionService(operationSectionRepo masteroperationrepository.OperationSectionRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationSectionService {
	return &OperationSectionServiceImpl{
		operationSectionRepo: operationSectionRepo,
		DB:                   db,
		RedisClient:          redisClient,
	}
}

func (s *OperationSectionServiceImpl) GetAllOperationSectionList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationSectionRepo.GetAllOperationSectionList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetSectionCodeByGroupId(GroupId int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationSectionRepo.GetSectionCodeByGroupId(tx, GroupId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionName(group_id int, section_code string) (masteroperationpayloads.OperationSectionNameResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationSectionRepo.GetOperationSectionName(tx, group_id, section_code)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) SaveOperationSection(req masteroperationpayloads.OperationSectionRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.OperationSectionId != 0 {
		_, err := s.operationSectionRepo.GetOperationSectionById(tx, req.OperationSectionId)
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation section not found",
			}
		}
	}

	if len(req.OperationSectionCode) > 3 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Operation Code max 3 characters",
		}
	}

	results, err := s.operationSectionRepo.SaveOperationSection(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save operation section",
		}
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionById(id int) (masteroperationpayloads.OperationSectionListResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationSectionRepo.GetOperationSectionById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) ChangeStatusOperationSection(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.operationSectionRepo.ChangeStatusOperationSection(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionDropDown(operationGroupId int) ([]masteroperationpayloads.OperationSectionDropDown, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.operationSectionRepo.GetOperationSectionDropDown(tx, operationGroupId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}
