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
	results, err := s.AgreementRepo.GetAgreementById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) SaveAgreement(req masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.AgreementId != 0 {
		_, err := s.AgreementRepo.GetAgreementById(tx, req.AgreementId)

		if err != nil {
			return masterentities.Agreement{}, err
		}
	}

	results, err := s.AgreementRepo.SaveAgreement(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return masterentities.Agreement{}, err
	}

	return results, nil
}

func (s *AgreementServiceImpl) UpdateAgreement(id int, req masterpayloads.AgreementRequest) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.AgreementRepo.GetAgreementById(tx, id)
	if err != nil {
		return masterentities.Agreement{}, err
	}

	results, err := s.AgreementRepo.UpdateAgreement(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.Agreement{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) ChangeStatusAgreement(Id int) (masterentities.Agreement, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.AgreementRepo.ChangeStatusAgreement(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.Agreement{}, err
	}
	return entity, nil
}

func (s *AgreementServiceImpl) GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllAgreement(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) AddDiscountGroup(id int, req masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.AddDiscountGroup(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementDiscountGroupDetail{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) UpdateDiscountGroup(id int, discountGroupId int, req masterpayloads.DiscountGroupRequest) (masterentities.AgreementDiscountGroupDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.UpdateDiscountGroup(tx, id, discountGroupId, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementDiscountGroupDetail{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) DeleteDiscountGroup(id int, discountGroupId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.AgreementRepo.DeleteDiscountGroup(tx, id, discountGroupId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) AddItemDiscount(id int, req masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.AddItemDiscount(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementItemDetail{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) UpdateItemDiscount(id int, itemDiscountId int, req masterpayloads.ItemDiscountRequest) (masterentities.AgreementItemDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.UpdateItemDiscount(tx, id, itemDiscountId, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementItemDetail{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) DeleteItemDiscount(id int, itemDiscountId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.AgreementRepo.DeleteItemDiscount(tx, id, itemDiscountId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) AddDiscountValue(id int, req masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.AddDiscountValue(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementDiscount{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) UpdateDiscountValue(id int, discountValueId int, req masterpayloads.DiscountValueRequest) (masterentities.AgreementDiscount, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.AgreementRepo.UpdateDiscountValue(tx, id, discountValueId, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.AgreementDiscount{}, err
	}
	return results, nil
}

func (s *AgreementServiceImpl) DeleteDiscountValue(id int, discountValueId int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.AgreementRepo.DeleteDiscountValue(tx, id, discountValueId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return err
	}
	return nil
}

func (s *AgreementServiceImpl) GetAllDiscountGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllDiscountGroup(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetAllItemDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllItemDiscount(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetAllDiscountValue(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.AgreementRepo.GetAllDiscountValue(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *AgreementServiceImpl) GetDiscountGroupAgreementById(agreementID, groupID int) (masterpayloads.DiscountGroupRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.AgreementRepo.GetDiscountGroupAgreementById(tx, agreementID, groupID)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterpayloads.DiscountGroupRequest{}, err
	}
	return result, nil
}

func (s *AgreementServiceImpl) GetDiscountItemAgreementById(agreementID, itemID int) (masterpayloads.ItemDiscountRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.AgreementRepo.GetDiscountItemAgreementById(tx, agreementID, itemID)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterpayloads.ItemDiscountRequest{}, err
	}
	return result, nil
}

func (s *AgreementServiceImpl) GetDiscountValueAgreementById(agreementID, valueID int) (masterpayloads.DiscountValueRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.AgreementRepo.GetDiscountValueAgreementById(tx, agreementID, valueID)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterpayloads.DiscountValueRequest{}, err
	}
	return result, nil
}
