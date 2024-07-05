package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"

	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemPackageDetailServiceImpl struct {
	ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository
	DB                    *gorm.DB
	RedisClient           *redis.Client // Redis client
}

func StartItemPackageDetailService(ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemPackageDetailService {
	return &ItemPackageDetailServiceImpl{
		ItemPackageDetailRepo: ItemPackageDetailRepo,
		DB:                    db,
		RedisClient:           redisClient,
	}
}

// ActivateItemPackageDetail implements masteritemservice.ItemPackageDetailService.
func (s *ItemPackageDetailServiceImpl) ActivateItemPackageDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.ActivateItemPackageDetail(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

// DeactiveItemPackageDetail implements masteritemservice.ItemPackageDetailService.
func (s *ItemPackageDetailServiceImpl) DeactiveItemPackageDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.DeactiveItemPackageDetail(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) ChangeStatusItemPackageDetail(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ItemPackageDetailRepo.GetItemPackageDetailById(tx, id)

	if err != nil {
		return false, err
	}

	results, err := s.ItemPackageDetailRepo.ChangeStatusItemPackageDetail(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailByItemPackageId(tx, itemPackageId, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailById(itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailById(tx, itemPackageDetailId)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) CreateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.CreateItemPackageDetailByItemPackageId(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) UpdateItemPackageDetail(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.UpdateItemPackageDetail(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}
