package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"

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

// GetItemLevelDropDown implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelDropDown(itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetItemLevelDropDown(tx, itemLevel)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.ItemLevelId != 0 {
		_, err := s.structItemLevelRepo.GetById(tx, request.ItemLevelId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.structItemLevelRepo.Save(tx, request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *ItemLevelServiceImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.structItemLevelRepo.GetAll(tx, request, pages)

	if err != nil {
		return get, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevelId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.structItemLevelRepo.GetById(tx, itemLevelId)

	if err != nil {
		return false, err
	}

	change_status, err := s.structItemLevelRepo.ChangeStatus(tx, itemLevelId)

	if err != nil {
		return change_status, err
	}

	return true, nil
}
