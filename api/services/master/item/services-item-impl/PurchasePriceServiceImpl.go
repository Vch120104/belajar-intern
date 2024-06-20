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

type PurchasePriceServiceImpl struct {
	PurchasePriceRepo masteritemrepository.PurchasePriceRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartPurchasePriceService(PurchasePriceRepo masteritemrepository.PurchasePriceRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.PurchasePriceService {
	return &PurchasePriceServiceImpl{
		PurchasePriceRepo: PurchasePriceRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePrice(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePrice(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) SavePurchasePrice(req masteritempayloads.PurchasePriceRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.PurchasePriceRepo.SavePurchasePrice(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceById(id int) (masteritempayloads.PurchasePriceRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.PurchasePriceRepo.GetPurchasePriceById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) AddPurchasePrice(req masteritempayloads.PurchasePriceDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.PurchasePriceRepo.AddPurchasePrice(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *PurchasePriceServiceImpl) GetAllPurchasePriceDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetAllPurchasePriceDetail(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *PurchasePriceServiceImpl) GetPurchasePriceDetailById(id int, pages pagination.Pagination) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.PurchasePriceRepo.GetPurchasePriceDetailById(tx, id, pages)
	if err != nil {
		return nil, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

// DeletePurchasePrice deletes an item location by ID
func (s *PurchasePriceServiceImpl) DeletePurchasePrice(id int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.PurchasePriceRepo.DeletePurchasePrice(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PurchasePriceServiceImpl) ChangeStatusPurchasePrice(Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Ubah status
	entity, err := s.PurchasePriceRepo.ChangeStatusPurchasePrice(tx, Id)
	if err != nil {
		return masteritementities.PurchasePrice{}, err
	}

	return entity, nil
}
