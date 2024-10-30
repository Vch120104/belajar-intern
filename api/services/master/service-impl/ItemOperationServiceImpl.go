package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemOperationServiceImpl struct {
	ItemOperationRepository masterrepository.ItemOperationRepository
	DB                      *gorm.DB
	RedisClient             *redis.Client // Redis client
}

func StartItemOperationService(ItemOperationRepo masterrepository.ItemOperationRepository, db *gorm.DB, redisClient *redis.Client) masterservice.ItemOperationService {
	return &ItemOperationServiceImpl{
		ItemOperationRepository: ItemOperationRepo,
		DB:                      db,
		RedisClient:             redisClient,
	}
}

func (s *ItemOperationServiceImpl) GetAllItemOperation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemOperationRepository.GetAllItemOperation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return pages, err
	}
	return result, nil
}

func (s *ItemOperationServiceImpl) GetByIdItemOperation(id int) (masterpayloads.ItemOperationPost, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemOperationRepository.GetByIdItemOperation(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterpayloads.ItemOperationPost{}, err
	}
	return result, nil
}

func (s *ItemOperationServiceImpl) PostItemOperation(req masterpayloads.ItemOperationPost) (masterentities.ItemOperation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemOperationRepository.PostItemOperation(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.ItemOperation{}, err
	}
	return result, nil
}

func (s *ItemOperationServiceImpl) DeleteItemOperation(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemOperationRepository.DeleteItemOperation(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *ItemOperationServiceImpl) UpdateItemOperation(id int, req masterpayloads.ItemOperationPost) (masterentities.ItemOperation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemOperationRepository.UpdateItemOperation(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.ItemOperation{}, err
	}
	return result, nil
}
