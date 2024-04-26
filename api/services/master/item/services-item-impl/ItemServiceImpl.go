package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemServiceImpl struct {
	itemRepo    masteritemrepository.ItemRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartItemService(itemRepo masteritemrepository.ItemRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemService {
	return &ItemServiceImpl{
		itemRepo:    itemRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition) []masteritempayloads.ItemLookup {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	// Definisikan parameter pagination
	paginate := pagination.Pagination{
		Limit: 10, // Sesuaikan dengan nilai limit yang diinginkan
		Page:  1,  // Sesuaikan dengan halaman yang diinginkan
	}

	// Panggil metode GetAllItem dari repository dengan argumen yang sesuai
	results, _, err := s.itemRepo.GetAllItem(tx, filterCondition, paginate)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemServiceImpl) GetAllItemLookup(queryParams map[string]string) []map[string]interface{} {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItemLookup(tx, queryParams)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemServiceImpl) GetItemById(Id int) masteritempayloads.ItemResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) []masteritempayloads.ItemResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}

func (s *ItemServiceImpl) GetItemCode(code string) []map[string]interface{} {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetItemCode(tx, code)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemId != 0 {
		_, err := s.itemRepo.GetItemById(tx, req.ItemId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.itemRepo.SaveItem(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.itemRepo.GetItemById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
