package masterserviceimpl

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
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
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	results, err := s.GmmDiscountSettingRepo.GetAllGmmDiscountSetting(tx)

	if err != nil {
		return results, err
	}
	return results, nil
}
