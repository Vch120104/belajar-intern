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

type ItemPackageServiceImpl struct {
	ItemPackageRepo masteritemrepository.ItemPackageRepository
	DB              *gorm.DB
	RedisClient     *redis.Client // Redis client
}

func StartItemPackageService(ItemPackageRepo masteritemrepository.ItemPackageRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemPackageService {
	return &ItemPackageServiceImpl{
		ItemPackageRepo: ItemPackageRepo,
		DB:              db,
		RedisClient:     redisClient,
	}
}

// GetItemPackageByCode implements masteritemservice.ItemPackageService.
func (s *ItemPackageServiceImpl) GetItemPackageByCode(itemPackageCode string) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemPackageRepo.GetItemPackageByCode(tx, itemPackageCode)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) GetAllItemPackage(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemPackageRepo.GetAllItemPackage(tx, internalFilterCondition, externalFilterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) GetItemPackageById(Id int) (masteritempayloads.GetItemPackageResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ItemPackageRepo.GetItemPackageById(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) SaveItemPackage(req masteritempayloads.SaveItemPackageRequest) (masteritementities.ItemPackage, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.ItemPackageRepo.SaveItemPackage(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageServiceImpl) ChangeStatusItemPackage(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.ItemPackageRepo.GetItemPackageById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.ItemPackageRepo.ChangeStatusItemPackage(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
