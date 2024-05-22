package masteroperationserviceimpl

import (
	"after-sales/api/exceptions"
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
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetAllOperationSectionList(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) GetSectionCodeByGroupId(GroupId int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetSectionCodeByGroupId(tx, GroupId)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *OperationSectionServiceImpl) GetOperationSectionName(group_id int, section_code string) (masteroperationpayloads.OperationSectionNameResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetOperationSectionName(tx, group_id, section_code)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *OperationSectionServiceImpl) SaveOperationSection(req masteroperationpayloads.OperationSectionRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

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
			StatusCode: http.StatusBadRequest,
			Message:    "Operation Section Code max 3 characters",
		}
	}

	results, err := s.operationSectionRepo.SaveOperationSection(tx, req)
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
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetOperationSectionById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationSectionServiceImpl) ChangeStatusOperationSection(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.ChangeStatusOperationSection(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}
