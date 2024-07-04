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

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstitute(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetByIdItemSubstitute(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstituteDetail(pages pagination.Pagination, id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetAllItemSubstituteDetail(tx, pages, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemSubstituteRepo.GetByIdItemSubstituteDetail(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.SaveItemSubstitute(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.SaveItemSubstituteDetail(tx, req, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstitute(filterCondition map[string]string, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse){
	tx:=s.Db.Begin()
	defer helper.CommitOrRollback(tx)
	result,totalPage,totalRows,err:=s.itemSubstituteRepo.GetAllItemSubstitute(tx,filterCondition,pages)
	if err != nil{
		return result,0,0,err
	}
	return result,totalPage,totalRows,nil
}

func (s *ItemSubstituteServiceImpl) ChangeStatusItemOperation(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.ChangeStatusItemOperation(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) DeactivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.DeactivateItemSubstituteDetail(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemSubstituteServiceImpl) ActivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.itemSubstituteRepo.ActivateItemSubstituteDetail(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}