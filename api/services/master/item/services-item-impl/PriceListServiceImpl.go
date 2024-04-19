package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

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

func (s *PriceListServiceImpl) GetPriceList(request masteritempayloads.PriceListGetAllRequest) []masteritempayloads.PriceListResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.priceListRepo.GetPriceList(tx, request)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *PriceListServiceImpl) GetPriceListById(Id int) masteritempayloads.PriceListResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.priceListRepo.GetPriceListById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *PriceListServiceImpl) SavePriceList(request masteritempayloads.PriceListResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.PriceListId != 0 {
		_, err := s.priceListRepo.GetPriceListById(tx, int(request.PriceListId))

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	result, err := s.priceListRepo.SavePriceList(tx, request)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}

func (s *PriceListServiceImpl) ChangeStatusPriceList(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.priceListRepo.GetPriceListById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	result, err := s.priceListRepo.ChangeStatusPriceList(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return result
}
