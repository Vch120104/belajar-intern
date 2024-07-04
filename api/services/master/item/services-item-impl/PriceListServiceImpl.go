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

type PriceListServiceImpl struct {
	priceListRepo masteritemrepository.PriceListRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartPriceListService(priceListRepo masteritemrepository.PriceListRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.PriceListService {
	return &PriceListServiceImpl{
		priceListRepo: priceListRepo,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *PriceListServiceImpl) GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.priceListRepo.GetPriceList(tx, request)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *PriceListServiceImpl) GetPriceListById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.priceListRepo.GetPriceListById(tx, Id)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *PriceListServiceImpl) SavePriceList(request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.PriceListId != 0 {
		_, err := s.priceListRepo.GetPriceListById(tx, int(request.PriceListId))

		if err != nil {
			return false, err
		}
	}

	result, err := s.priceListRepo.SavePriceList(tx, request)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PriceListServiceImpl) ChangeStatusPriceList(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.priceListRepo.GetPriceListById(tx, Id)

	if err != nil {
		return false, err
	}

	result, err := s.priceListRepo.ChangeStatusPriceList(tx, Id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PriceListServiceImpl) GetAllPriceListNew(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, total_page, total_rows, err := s.priceListRepo.GetAllPriceListNew(tx, filterCondition, pages)
	if err != nil {
		return nil, 0, 0, err
	}
	defer helper.CommitOrRollback(tx, err)

	return result, total_page, total_rows, nil
}

func (s *PriceListServiceImpl) DeactivatePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.priceListRepo.DeactivatePriceList(tx, id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)

	return result, nil
}

func (s *PriceListServiceImpl) ActivatePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.priceListRepo.ActivatePriceList(tx, id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)

	return result, nil
}

func (s *PriceListServiceImpl) DeletePriceList(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.priceListRepo.DeletePriceList(tx, id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)

	return result, nil
}
