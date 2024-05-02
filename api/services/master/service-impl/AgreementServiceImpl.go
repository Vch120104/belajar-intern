package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AgreementServiceImpl struct {
	AgreementRepo masterrepository.AgreementRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartAgreementService(AgreementRepo masterrepository.AgreementRepository, db *gorm.DB, redisClient *redis.Client) masterservice.AgreementService {
	return &AgreementServiceImpl{
		AgreementRepo: AgreementRepo,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *AgreementServiceImpl) GetAgreementById(id int) (masterpayloads.AgreementResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.AgreementRepo.GetAgreementById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) SaveAgreement(req masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.AgreementId != 0 {
		_, err := s.AgreementRepo.GetAgreementById(tx, req.AgreementId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.AgreementRepo.SaveAgreement(tx, req)

	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) ChangeStatusAgreement(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.AgreementRepo.GetAgreementById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.AgreementRepo.ChangeStatusAgreement(tx, Id)
	if err != nil {
		return results, nil
	}
	return true, nil
}

func (s *AgreementServiceImpl) GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllAgreement(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}
