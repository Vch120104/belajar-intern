package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemClassServiceImpl struct {
	itemRepo    masteritemrepository.ItemClassRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartItemClassService(itemRepo masteritemrepository.ItemClassRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemClassService {
	return &ItemClassServiceImpl{
		itemRepo:    itemRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

func (s *ItemClassServiceImpl) GetAllItemClass(filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItemClass(tx, filterCondition)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *ItemClassServiceImpl) GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemClassById(tx, Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemClassServiceImpl) SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemClassId != 0 {
		_, err := s.itemRepo.GetItemClassById(tx, req.ItemClassId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.itemRepo.SaveItemClass(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemClassServiceImpl) ChangeStatusItemClass(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.itemRepo.GetItemClassById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItemClass(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
