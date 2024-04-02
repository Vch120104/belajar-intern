package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptionsss_test "after-sales/api/expectionsss"
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

func (r *OperationCodeRepositoryImpl) GetAllOperationCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []masteroperationentities.OperationCode{}
	// var payloads []masteroperationpayloads.OperationCodeGetAll
	baseModelQuery := tx.Model(&entities)
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()
	if len(entities) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *OperationCodeRepositoryImpl) GetOperationCodeById(tx *gorm.DB, id int) (masteroperationpayloads.OperationCodeResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}
	rows, err := tx.Model(&entities).Where(masteroperationentities.OperationCode{OperationId: id}).First(&response).Rows()
	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()
	return response, nil
}

func (r *OperationCodeRepositoryImpl) SaveOperationCode(tx *gorm.DB, req masteroperationpayloads.OperationCodeSave) (bool, *exceptionsss_test.BaseErrorResponse) {
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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	if len(req.OperationCode) > 10 || len(req.OperationCode)>200 {
		// errMessage := "Operation Group Code max 2 characters"

		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,

			Err: errors.New(utils.BadRequestError),
		}
	}
	return true, nil
}

func (r *OperationCodeRepositoryImpl) ChangeStatusItemSubstitute(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteroperationentities.OperationCode

	result := tx.Model(&entities).
		Where("operation_id = ?", id).
		First(&entities)

		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
