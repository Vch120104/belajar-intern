package masteroperationserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

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

func (s *OperationSectionServiceImpl) GetAllOperationSectionList(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetAllOperationSectionList(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationSectionServiceImpl) GetSectionCodeByGroupId(GroupId int) []masteroperationpayloads.OperationSectionCodeResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetSectionCodeByGroupId(tx, GroupId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return results
}

func (s *OperationSectionServiceImpl) GetOperationSectionName(group_id int, section_code string) masteroperationpayloads.OperationSectionNameResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetOperationSectionName(tx, group_id, section_code)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return results
}

func (s *OperationSectionServiceImpl) SaveOperationSection(req masteroperationpayloads.OperationSectionRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.OperationSectionId != 0 {
		_, err := s.operationSectionRepo.GetOperationSectionById(tx, req.OperationSectionId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	if len(req.OperationSectionCode) > 3 {
		panic(exceptions.NewBadRequestError("Operation Section Code max 3 characters"))
	}
	results, err := s.operationSectionRepo.SaveOperationSection(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationSectionServiceImpl) GetOperationSectionById(id int) masteroperationpayloads.OperationSectionListResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.GetOperationSectionById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *OperationSectionServiceImpl) ChangeStatusOperationSection(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.operationSectionRepo.ChangeStatusOperationSection(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
