package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
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

// GetItemClassDropDownbyGroupId implements masteritemservice.ItemClassService.
func (s *ItemClassServiceImpl) GetItemClassDropDownbyGroupId(groupId int) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemClassDropDownbyGroupId(tx, groupId)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

// GetItemClassByCode implements masteritemservice.ItemClassService.
func (s *ItemClassServiceImpl) GetItemClassByCode(itemClassCode string) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemClassByCode(tx, itemClassCode)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

// GetItemClassDropDown implements masteritemservice.ItemClassService.
func (s *ItemClassServiceImpl) GetItemClassDropDown() ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetItemClassDropDown(tx)
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemClassServiceImpl) GetAllItemClass(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.itemRepo.GetAllItemClass(tx, filterCondition, pages)
	if err != nil {
		return nil, 0, 0, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, totalPages, totalRows, nil
}

func (s *ItemClassServiceImpl) GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemClassById(tx, Id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemClassServiceImpl) SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

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
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemClassServiceImpl) ChangeStatusItemClass(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.itemRepo.GetItemClassById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItemClass(tx, Id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}
