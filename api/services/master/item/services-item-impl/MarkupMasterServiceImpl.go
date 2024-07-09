package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MarkupMasterServiceImpl struct {
	markupRepo  masteritemrepository.MarkupMasterRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartMarkupMasterService(markupRepo masteritemrepository.MarkupMasterRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.MarkupMasterService {
	return &MarkupMasterServiceImpl{
		markupRepo:  markupRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.markupRepo.GetMarkupMasterList(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterById(id int) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.markupRepo.GetMarkupMasterById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) GetAllMarkupMasterIsActive() ([]masteritempayloads.MarkupMasterDropDownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.markupRepo.GetAllMarkupMasterIsActive(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) SaveMarkupMaster(req masteritempayloads.MarkupMasterResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.MarkupMasterId != 0 {
		_, err := s.markupRepo.GetMarkupMasterById(tx, req.MarkupMasterId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.markupRepo.SaveMarkupMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}
func (s *MarkupMasterServiceImpl) ChangeStatusMasterMarkupMaster(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.markupRepo.GetMarkupMasterById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.markupRepo.ChangeStatusMasterMarkupMaster(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return true, nil
}
func (s *MarkupMasterServiceImpl) GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.markupRepo.GetMarkupMasterByCode(tx, markupCode)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil

}
