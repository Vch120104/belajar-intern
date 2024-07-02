package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemSubstituteServiceImpl struct {
	itemSubstituteRepo masteritemrepository.ItemSubstituteRepository
	Db                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func StartItemSubstituteService(itemSubstituteRepo masteritemrepository.ItemSubstituteRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemSubstituteService {
	return &ItemSubstituteServiceImpl{
		itemSubstituteRepo: itemSubstituteRepo,
		Db:                 db,
		RedisClient:        redisClient,
	}
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination, from time.Time, to time.Time) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	results, err := s.itemSubstituteRepo.GetAllItemSubstitute(tx, filterCondition, pages, from, to)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstitute(id int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	result, err := s.itemSubstituteRepo.GetByIdItemSubstitute(tx, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstituteDetail(pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	result, err := s.itemSubstituteRepo.GetAllItemSubstituteDetail(tx, pages, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	result, err := s.itemSubstituteRepo.GetByIdItemSubstituteDetail(tx, id)

	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()

	result, err := s.itemSubstituteRepo.SaveItemSubstitute(tx, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()

	result, err := s.itemSubstituteRepo.SaveItemSubstituteDetail(tx, req, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) ChangeStatusItemSubstitute(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()

	result, err := s.itemSubstituteRepo.ChangeStatusItemSubstitute(tx, id)

	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) DeactivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()

	result, err := s.itemSubstituteRepo.DeactivateItemSubstituteDetail(tx, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemSubstituteServiceImpl) ActivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()

	result, err := s.itemSubstituteRepo.ActivateItemSubstituteDetail(tx, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}
