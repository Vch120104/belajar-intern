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

type OperationGroupRepositoryImpl struct {
}

func StartOperationGroupRepositoryImpl() masteroperationrepository.OperationGroupRepository {
	return &OperationGroupRepositoryImpl{}
}

func (r *OperationGroupRepositoryImpl) GetAllOperationGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	// entities := []masteroperationentities.OperationGroup{}
	// //define base model
	// baseModelQuery := tx.Model(&entities)
	// //apply where query
	// whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	// //apply pagination and execute
	// rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	// if len(entities) == 0 {
	// 	return pages, gorm.ErrRecordNotFound
	// }

	// if err != nil {
	// 	return pages, err
	// }

	// defer rows.Close()

	// pages.Rows = entities

	// return pages, nil
	entities := []masteroperationentities.OperationGroup{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *OperationGroupRepositoryImpl) GetOperationGroupById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupId: Id,
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

func (r *OperationGroupRepositoryImpl) GetOperationGroupByCode(tx *gorm.DB, Code string) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupCode: Code,
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

func (r *OperationGroupRepositoryImpl) SaveOperationGroup(tx *gorm.DB, req masteroperationpayloads.OperationGroupResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{
		IsActive:                  req.IsActive,
		OperationGroupId:          req.OperationGroupId,
		OperationGroupCode:        req.OperationGroupCode,
		OperationGroupDescription: req.OperationGroupDescription,
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

func (r *OperationGroupRepositoryImpl) ChangeStatusOperationGroup(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationGroup

	result := tx.Model(&entities).
		Where("operation_group_id = ?", Id).
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

func (r *OperationGroupRepositoryImpl) GetOperationGroupDropDown(tx *gorm.DB) ([]masteroperationpayloads.OperationGroupDropDownResponse, *exceptions.BaseErrorResponse) {

	var operationGroupDropDownResponse []masteroperationpayloads.OperationGroupDropDownResponse

	err := tx.Model(&masteroperationentities.OperationGroup{}).
		Select("operation_group_id", "CONCAT(operation_group_code, ' - ', operation_group_description) as operation_group_code").
		Find(&operationGroupDropDownResponse)
	if err.Error != nil {
		return operationGroupDropDownResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error,
		}
	}
	return operationGroupDropDownResponse, nil
}
