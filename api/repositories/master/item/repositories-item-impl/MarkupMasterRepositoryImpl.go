package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type MarkupMasterRepositoryImpl struct {
}

func StartMarkupMasterRepositoryImpl() masteritemrepository.MarkupMasterRepository {
	return &MarkupMasterRepositoryImpl{}
}

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritementities.MarkupMaster

	baseModelQuery := tx.Model(&responses)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&responses, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *MarkupMasterRepositoryImpl) GetMarkupMasterById(tx *gorm.DB, Id int) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.MarkupMaster{}
	response := masteritempayloads.MarkupMasterResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.MarkupMaster{
			MarkupMasterId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *MarkupMasterRepositoryImpl) GetAllMarkupMasterIsActive(tx *gorm.DB) ([]masteritempayloads.MarkupMasterDropDownResponse, *exceptions.BaseErrorResponse) {
	MarkupMasters := masteritementities.MarkupMaster{}
	response := []masteritempayloads.MarkupMasterDropDownResponse{}

	err := tx.Model(&MarkupMasters).Select("mtr_markup_master.*, CONCAT(markup_code, ' - ', markup_description) AS markup_code_description").Where("is_active = 'true'").Find(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return response, nil
}

func (r *MarkupMasterRepositoryImpl) SaveMarkupMaster(tx *gorm.DB, req masteritempayloads.MarkupMasterResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.MarkupMaster{
		IsActive:          req.IsActive,
		MarkupMasterId:    req.MarkupMasterId,
		MarkupCode:        req.MarkupCode,
		MarkupDescription: req.MarkupDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}
func (r *MarkupMasterRepositoryImpl) ChangeStatusMasterMarkupMaster(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.MarkupMaster

	result := tx.Model(&entities).
		Where("markup_master_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
func (r *MarkupMasterRepositoryImpl) GetMarkupMasterByCode(tx *gorm.DB, markupCode string) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.MarkupMasterResponse{}
	var entities masteritementities.MarkupMaster
	rows, err := tx.Model(&entities).
		Where("markup_code = ?", markupCode).
		First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}
