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

type ItemPriceCodeServiceImpl struct {
	ItemPriceCodeRepo masteritemrepository.ItemPriceCodeRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartItemPriceCodeService(ItemPriceCodeRepo masteritemrepository.ItemPriceCodeRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemPriceCodeService {
	return &ItemPriceCodeServiceImpl{
		ItemPriceCodeRepo: ItemPriceCodeRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *ItemPriceCodeServiceImpl) GetAllItemPriceCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.ItemPriceCodeRepo.GetAllItemPriceCode(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemPriceCodeServiceImpl) GetByIdItemPriceCode(id int) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.GetByIdItemPriceCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemPriceCodeServiceImpl) GetByCodeItemPriceCode(itemPriceCode string) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.GetByCodeItemPriceCode(tx, itemPriceCode)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemPriceCodeServiceImpl) SaveItemPriceCode(req masteritempayloads.SaveItemPriceCode) (masteritementities.ItemPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.SaveItemPriceCode(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemPriceCodeServiceImpl) DeleteItemPriceCode(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.DeleteItemPriceCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *ItemPriceCodeServiceImpl) UpdateItemPriceCode(itemPriceId int, req masteritempayloads.UpdateItemPriceCode) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.UpdateItemPriceCode(tx, itemPriceId, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *ItemPriceCodeServiceImpl) ChangeStatusItemPriceCode(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.ItemPriceCodeRepo.ChangeStatusItemPriceCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}
