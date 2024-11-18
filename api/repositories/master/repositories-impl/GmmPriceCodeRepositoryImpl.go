package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type GmmPriceCodeRepositoryImpl struct {
}

func StartGmmPriceCodeRepositoryImpl() masterrepository.GmmPriceCodeRepository {
	return &GmmPriceCodeRepositoryImpl{}
}

func (r *GmmPriceCodeRepositoryImpl) GetAllGmmPriceCode(tx *gorm.DB) ([]masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}
	response := []masterpayloads.GmmPriceCodeResponse{}

	err := tx.Model(&entities).Scan(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	return response, nil
}

func (r *GmmPriceCodeRepositoryImpl) GetGmmPriceCodeById(tx *gorm.DB, id int) (masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}
	response := masterpayloads.GmmPriceCodeResponse{}

	err := tx.Model(&entities).Where(masterentities.GmmPriceCode{GmmPriceCodeId: id}).First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	return response, nil
}

func (r *GmmPriceCodeRepositoryImpl) GetGmmPriceCodeDropdown(tx *gorm.DB) ([]masterpayloads.GmmPriceCodeDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}
	payloads := []masterpayloads.GmmPriceCodeResponse{}
	response := []masterpayloads.GmmPriceCodeDropdownResponse{}

	err := tx.Model(&entities).Where(masterentities.GmmPriceCode{IsActive: true}).Scan(&payloads).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	for _, data := range payloads {
		resp := masterpayloads.GmmPriceCodeDropdownResponse{
			IsActive:         data.IsActive,
			GmmPriceCodeId:   data.GmmPriceCodeId,
			GmmPriceCode:     data.GmmPriceCode,
			GmmPriceDesc:     data.GmmPriceDesc,
			GmmPriceCodeDesc: data.GmmPriceCode + " - " + data.GmmPriceDesc,
		}
		response = append(response, resp)
	}

	return response, nil
}

func (r *GmmPriceCodeRepositoryImpl) SaveGmmPriceCode(tx *gorm.DB, req masterpayloads.GmmPriceCodeSaveRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{
		IsActive:     true,
		GmmPriceCode: req.GmmPriceCode,
		GmmPriceDesc: req.GmmPriceDesc,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "gmm price code already exist",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving gmm price code",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *GmmPriceCodeRepositoryImpl) UpdateGmmPriceCode(tx *gorm.DB, id int, req masterpayloads.GmmPriceCodeUpdateRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}

	err := tx.Model(&entities).Where(masterentities.GmmPriceCode{GmmPriceCodeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	entities.GmmPriceCode = req.GmmPriceCode
	entities.GmmPriceDesc = req.GmmPriceDesc

	err = tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "gmm price code already exist",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving gmm price code",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *GmmPriceCodeRepositoryImpl) ChangeStatusGmmPriceCode(tx *gorm.DB, id int) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}

	err := tx.Model(&entities).Where(masterentities.GmmPriceCode{GmmPriceCodeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	entities.IsActive = !entities.IsActive

	err = tx.Save(&entities).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "gmm price code already exist",
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error saving gmm price code",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *GmmPriceCodeRepositoryImpl) DeleteGmmPriceCode(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	entities := masterentities.GmmPriceCode{}

	err := tx.Model(&entities).Where(masterentities.GmmPriceCode{GmmPriceCodeId: id}).First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching gmm price code",
			Err:        err,
		}
	}

	err = tx.Delete(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error deleting gmm price code",
			Err:        err,
		}
	}

	return true, nil
}
