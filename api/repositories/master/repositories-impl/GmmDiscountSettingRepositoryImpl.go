package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"net/http"

	"gorm.io/gorm"
)

type GmmDiscountSettingRepositoryImpl struct {
}

func StartGmmDiscountSettingRepositoryImpl() masterrepository.GmmDiscountSettingRepository {
	return &GmmDiscountSettingRepositoryImpl{}
}

func (r *GmmDiscountSettingRepositoryImpl) GetAllGmmDiscountSetting(tx *gorm.DB) ([]masterpayloads.GmmDiscountSettingResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmDiscountSetting{}
	response := []masterpayloads.GmmDiscountSettingResponse{}

	err := tx.Model(&entities).Scan(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}
