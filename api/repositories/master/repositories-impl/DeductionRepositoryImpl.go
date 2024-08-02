package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type DeductionRepositoryImpl struct {
}

func StartDeductionRepositoryImpl() masterrepository.DeductionRepository {
	return &DeductionRepositoryImpl{}
}

func (r *DeductionRepositoryImpl) GetAllDeduction(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masterentities.DeductionList{}
	response := []masterpayloads.DeductionListResponse{}

	baseModelQuery := tx.Model(&entities).Scan(&response)

	wherequery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, wherequery)).Scan(&response).Rows()

	if len(response) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = response

	return pages, nil
}

func (r *DeductionRepositoryImpl) GetDeductionById(tx *gorm.DB, Id int) (masterpayloads.DeductionListResponse, *exceptions.BaseErrorResponse) {

	entities := masterentities.DeductionList{}

	response := masterpayloads.DeductionListResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.DeductionList{
			DeductionId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *DeductionRepositoryImpl) GetAllDeductionDetail(tx *gorm.DB, pages pagination.Pagination, Id int) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masterentities.DeductionDetail{}

	response := []masterpayloads.DeductionDetailResponse{}

	baseModelQuery := tx.Model(&entities).
		Where(masterentities.DeductionDetail{
			DeductionId: Id},
		)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&response).Rows()

	if len(response) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = response

	return pages, nil
}

func (r *DeductionRepositoryImpl) GetByIdDeductionDetail(tx *gorm.DB, Id int) (masterpayloads.DeductionDetailResponse, *exceptions.BaseErrorResponse) {

	entities := masterentities.DeductionDetail{}

	response := masterpayloads.DeductionDetailResponse{}

	rows, err := tx.Model(&entities).Where(
		masterentities.DeductionDetail{
			DeductionDetailId: Id,
		}).First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	return response, nil
}

func (r *DeductionRepositoryImpl) SaveDeductionList(tx *gorm.DB, request masterpayloads.DeductionListResponse) (masterentities.DeductionList, *exceptions.BaseErrorResponse) {

	entities := masterentities.DeductionList{
		DeductionCode: request.DeductionCode,
		DeductionName: request.DeductionName,
		EffectiveDate: request.EffectiveDate,
	}

	err := tx.Where(entities).Assign(entities).FirstOrCreate(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return masterentities.DeductionList{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return masterentities.DeductionList{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities, nil
}

func (r *DeductionRepositoryImpl) SaveDeductionDetail(tx *gorm.DB, request masterpayloads.DeductionDetailResponse) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse) {
	condition := masterentities.DeductionDetail{
		DeductionId:          request.DeductionId,
		DeductionDetailLevel: request.DeductionDetailLevel,
	}

	entities := masterentities.DeductionDetail{
		DeductionId:          request.DeductionId,
		DeductionDetailLevel: request.DeductionDetailLevel,
		LimitDays:            request.LimitDays,
		DeductionPercent:     request.DeductionPercent,
	}

	err := tx.Where(entities).Assign(condition).FirstOrCreate(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return masterentities.DeductionDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return masterentities.DeductionDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	return entities, nil
}

func (r *DeductionRepositoryImpl) ChangeStatusDeduction(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var entities masterentities.DeductionList
	result := tx.Model(&entities).
		Where(masterentities.DeductionList{DeductionId: Id}).
		First(&entities)

	if result.Error != nil {
		return nil, &exceptions.BaseErrorResponse{
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
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}
	results := map[string]interface{}{
		"is_active":    entities.IsActive,
		"deduction_id": entities.DeductionId,
	}

	return results, nil
}

func (r *DeductionRepositoryImpl) UpdateDeductionDetail(tx *gorm.DB, id int, req masterpayloads.DeductionDetailUpdate) (masterentities.DeductionDetail, *exceptions.BaseErrorResponse) {
	var entities masterentities.DeductionDetail
	err := tx.Model(&entities).Where("deduction_detail_id = ?", id).Updates(map[string]interface{}{
		"limit_days":        req.LimitDays,
		"deduction_percent": req.DeductionPercent,
	}).Error
	if err != nil {
		return masterentities.DeductionDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	err = tx.Where("deduction_detail_id = ?", id).First(&entities).Error
	if err != nil{
		return masterentities.DeductionDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return entities, nil
}
