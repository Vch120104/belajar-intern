package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
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

func (s *AgreementServiceImpl) GetAgreementById(id int) (masterpayloads.AgreementRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.AgreementRepo.GetAgreementById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) SaveAgreement(req masterpayloads.AgreementRequest) (bool, *exceptions.BaseErrorResponse) {
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

func (s *AgreementServiceImpl) ChangeStatusAgreement(Id int) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Ubah status
	entity, err := s.AgreementRepo.ChangeStatusAgreement(tx, Id)
	if err != nil {
		return masterentities.Agreement{}, err
	}

	return entity, nil
}

func (s *AgreementServiceImpl) GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllAgreement(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) AddDiscountGroup(id int, req masterpayloads.DiscountGroupRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.AddDiscountGroup(tx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) DeleteDiscountGroup(id int, discountGroupId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.DeleteDiscountGroup(tx, id, discountGroupId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) AddItemDiscount(id int, req masterpayloads.ItemDiscountRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.AddItemDiscount(tx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) DeleteItemDiscount(id int, itemDiscountId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.DeleteItemDiscount(tx, id, itemDiscountId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) AddDiscountValue(id int, req masterpayloads.DiscountValueRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.AddDiscountValue(tx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) DeleteDiscountValue(id int, discountValueId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.AgreementRepo.DeleteDiscountValue(tx, id, discountValueId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) GetAllDiscountGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllDiscountGroup(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetAllItemDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllItemDiscount(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetAllDiscountValue(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllDiscountValue(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetDiscountGroupAgreementById(agreementID, groupID int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.AgreementRepo.GetDiscountGroupAgreementById(tx, agreementID, groupID)
	if err != nil {
		return masterpayloads.DiscountGroupRequest{}, err
	}
	return result, nil
}

func (s *AgreementServiceImpl) GetDiscountItemAgreementById(agreementID, itemID int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.AgreementRepo.GetDiscountItemAgreementById(tx, agreementID, itemID)
	if err != nil {
		return masterpayloads.ItemDiscountRequest{}, err
	}
	return result, nil
}

func (s *AgreementServiceImpl) GetDiscountValueAgreementById(agreementID, valueID int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.AgreementRepo.GetDiscountValueAgreementById(tx, agreementID, valueID)
	if err != nil {
		return masterpayloads.DiscountValueRequest{}, err
	}
	return result, nil
}
