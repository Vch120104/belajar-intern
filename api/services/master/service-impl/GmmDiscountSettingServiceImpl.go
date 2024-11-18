package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type GmmDiscountSettingServiceImpl struct {
	GmmDiscountSettingRepo masterrepository.GmmDiscountSettingRepository
	DB                     *gorm.DB
}

func StartGmmDiscountSettingServiceImpl(gmmDiscountSettingRepo masterrepository.GmmDiscountSettingRepository, db *gorm.DB) masterservice.GmmDiscountSettingService {
	return &GmmDiscountSettingServiceImpl{
		GmmDiscountSettingRepo: gmmDiscountSettingRepo,
		DB:                     db,
	}
}

func (s *GmmDiscountSettingServiceImpl) GetAllGmmDiscountSetting() ([]masterpayloads.GmmDiscountSettingResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.GmmDiscountSettingRepo.GetAllGmmDiscountSetting(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}
