package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type OperationCodeRepositoryImpl struct {
}

func StartOperationCodeRepositoryImpl() masteroperationrepository.OperationCodeRepository {
	return &OperationCodeRepositoryImpl{}
}

func (r *OperationCodeRepositoryImpl) GetAllOperationCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteroperationentities.OperationCode{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()
	if len(entities) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (*OperationCodeRepositoryImpl) GetAllOperationCodeDropDown(tx *gorm.DB) ([]masteroperationpayloads.OperationCodeGetAll, *exceptions.BaseErrorResponse) {
	baseModelQuery := tx.Model(&masteroperationentities.OperationCode{}).Select("operation_id, operation_code, operation_name, is_active")
	rows, err := baseModelQuery.Rows()
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	var responses []masteroperationpayloads.OperationCodeGetAll
	for rows.Next() {
		var operationId int
		var operationCode, operationName string
		var isActive bool

		err := rows.Scan(&operationId, &operationCode, &operationName, &isActive)
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		responseMap := masteroperationpayloads.OperationCodeGetAll{
			OperationCode: operationCode,
			OperationName: operationName,
			IsActive:      isActive,
		}
		responses = append(responses, responseMap)
	}

	return responses, nil
}

func (r *OperationCodeRepositoryImpl) GetOperationCodeById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationCode{
			OperationId: Id,
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

func (r *OperationCodeRepositoryImpl) GetOperationCodeByCode(tx *gorm.DB, code string) (masteroperationpayloads.OperationCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}

	err := tx.Model(&entities).
		Where(masteroperationentities.OperationCode{
			OperationCode: code,
		}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *OperationCodeRepositoryImpl) SaveOperationCode(tx *gorm.DB, req masteroperationpayloads.OperationCodeSave) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationCode{
		IsActive:                req.IsActive,
		OperationCode:           req.OperationCode,
		OperationName:           req.OperationName,
		OperationUsingIncentive: req.OperationUsingIncentive,
		OperationUsingActual:    req.OperationUsingActual,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return masteroperationentities.OperationCode{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {
			return masteroperationentities.OperationCode{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	if len(req.OperationCode) > 10 || len(req.OperationCode) > 200 {
		// errMessage := "Operation Group Code max 2 characters"

		return masteroperationentities.OperationCode{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,

			Err: errors.New(utils.BadRequestError),
		}
	}
	return entities, nil
}

func (r *OperationCodeRepositoryImpl) ChangeStatusItemCode(tx *gorm.DB, id int) (masteroperationentities.OperationCode, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationCode

	result := tx.Model(&entities).
		Where("operation_id = ?", id).
		First(&entities)

	if result.Error != nil {
		return masteroperationentities.OperationCode{}, &exceptions.BaseErrorResponse{
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
		return masteroperationentities.OperationCode{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *OperationCodeRepositoryImpl) UpdateItemCode(tx *gorm.DB, id int, req masteroperationpayloads.OperationCodeUpdate) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationCode

	err := tx.Model(&entities).Where("operation_id = ?", id).Updates(map[string]interface{}{
		"operation_name":            req.OperationName,
		"operation_using_actual":    req.OperationUsingActual,
		"operation_using_incentive": req.OperationUsingIncentive,
	}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	return true, nil
}
