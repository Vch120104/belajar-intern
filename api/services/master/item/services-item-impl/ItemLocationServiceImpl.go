package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
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

type ItemLocationServiceImpl struct {
	ItemLocationRepo masteritemrepository.ItemLocationRepository
	DB               *gorm.DB
	RedisClient      *redis.Client // Redis client
}

func StartItemLocationService(ItemLocationRepo masteritemrepository.ItemLocationRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemLocationService {
	return &ItemLocationServiceImpl{
		ItemLocationRepo: ItemLocationRepo,
		DB:               db,
		RedisClient:      redisClient,
	}
}

func (s *ItemLocationServiceImpl) GetAllItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.GetAllItemLocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemLocationServiceImpl) SaveItemLocation(req masteritempayloads.ItemLocationRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemLocationRepo.SaveItemLocation(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemLocationServiceImpl) AddItemLocation(id int, req masteritempayloads.ItemLocationDetailRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.ItemLocationRepo.AddItemLocation(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemLocationServiceImpl) GetItemLocationById(id int) (masteritempayloads.ItemLocationRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemLocationRepo.GetItemLocationById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemLocationServiceImpl) GetAllItemLocationDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.GetAllItemLocationDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemLocationServiceImpl) PopupItemLocation(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemLocationRepo.PopupItemLocation(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// DeleteItemLocation deletes an item location by ID
func (s *ItemLocationServiceImpl) DeleteItemLocation(id int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.ItemLocationRepo.DeleteItemLocation(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemLocationServiceImpl) GetAllItemLoc(filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, totalpages, totalrows, err := s.ItemLocationRepo.GetAllItemLoc(tx, filtercondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, 0, 0, err
	}
	return result, totalpages, totalrows, nil
}

func (s *ItemLocationServiceImpl) GetByIdItemLoc(id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.GetByIdItemLoc(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemLocationServiceImpl) SaveItemLoc(req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.SaveItemLoc(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.ItemLocation{}, err
	}
	return result, nil
}

func (s *ItemLocationServiceImpl) DeleteItemLoc(ids []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemLocationRepo.DeleteItemLoc(tx, ids)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}
