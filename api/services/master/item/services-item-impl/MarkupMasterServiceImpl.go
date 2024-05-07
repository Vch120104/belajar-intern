package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
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

func (s *MarkupMasterServiceImpl) GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupMasterList(tx, filter, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterById(id int) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupMasterById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) SaveMarkupMaster(req masteritempayloads.MarkupMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.MarkupMasterId != 0 {
		_, err := s.markupRepo.GetMarkupMasterById(tx, req.MarkupMasterId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.markupRepo.SaveMarkupMaster(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}
func (s *MarkupMasterServiceImpl) ChangeStatusMasterMarkupMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.markupRepo.GetMarkupMasterById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.markupRepo.ChangeStatusMasterMarkupMaster(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
func (s *MarkupMasterServiceImpl) GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.markupRepo.GetMarkupMasterByCode(tx, markupCode)
	if err != nil {
		return result, err
	}
	return result, nil

}
