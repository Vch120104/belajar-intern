package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type GmmPriceCodeServiceImpl struct {
	GmmPriceCodeRepo masterrepository.GmmPriceCodeRepository
	DB               *gorm.DB
}

func StartGmmPriceCodeServiceImpl(gmmPriceCodeRepo masterrepository.GmmPriceCodeRepository, db *gorm.DB) masterservice.GmmPriceCodeService {
	return &GmmPriceCodeServiceImpl{
		GmmPriceCodeRepo: gmmPriceCodeRepo,
		DB:               db,
	}
}

func (s *GmmPriceCodeServiceImpl) GetAllGmmPriceCode() ([]masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.GetAllGmmPriceCode(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) GetGmmPriceCodeById(id int) (masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.GetGmmPriceCodeById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) GetGmmPriceCodeDropdown() ([]masterpayloads.GmmPriceCodeDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.GetGmmPriceCodeDropdown(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) SaveGmmPriceCode(req masterpayloads.GmmPriceCodeSaveRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.SaveGmmPriceCode(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) UpdateGmmPriceCode(id int, req masterpayloads.GmmPriceCodeUpdateRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.UpdateGmmPriceCode(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) ChangeStatusGmmPriceCode(id int) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.ChangeStatusGmmPriceCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) DeleteGmmPriceCode(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmPriceCodeRepo.DeleteGmmPriceCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}
