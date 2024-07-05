package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemLevelServiceImpl struct {
	structItemLevelRepo masteritemlevelrepo.ItemLevelRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func StartItemLevelService(itemlevelrepo masteritemlevelrepo.ItemLevelRepository, db *gorm.DB, redisClient *redis.Client) masteritemlevelservice.ItemLevelService {
	return &ItemLevelServiceImpl{
		structItemLevelRepo: itemlevelrepo,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

// GetItemLevelLookUp implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.structItemLevelRepo.GetItemLevelLookUp(tx, filter, pages, itemClassId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

// GetItemLevelDropDown implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.structItemLevelRepo.GetItemLevelDropDown(tx, itemLevel)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.ItemLevelId != 0 {
		_, err := s.structItemLevelRepo.GetById(tx, request.ItemLevelId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.structItemLevelRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *ItemLevelServiceImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.structItemLevelRepo.GetById(tx, itemLevelId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.structItemLevelRepo.GetAll(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevelId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		return false, err
	}

	change_status, err := s.structItemLevelRepo.ChangeStatus(tx, itemLevelId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return change_status, err
	}

	return true, nil
}
