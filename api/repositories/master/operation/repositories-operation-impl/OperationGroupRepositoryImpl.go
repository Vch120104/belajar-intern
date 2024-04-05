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

type OperationGroupRepositoryImpl struct {
}

func StartOperationGroupRepositoryImpl() masteroperationrepository.OperationGroupRepository {
	return &OperationGroupRepositoryImpl{}
}

func (*OperationGroupRepositoryImpl) GetAllOperationGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
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
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}
	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (*OperationGroupRepositoryImpl) GetAllOperationGroupIsActive(tx *gorm.DB) ([]masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	var OperationGroups []masteroperationentities.OperationGroup
	response := []masteroperationpayloads.OperationGroupResponse{}

	err := tx.Model(&OperationGroups).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (*OperationGroupRepositoryImpl) GetOperationGroupById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (*OperationGroupRepositoryImpl) GetOperationGroupByCode(tx *gorm.DB, Code string) (masteroperationpayloads.OperationGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (*OperationGroupRepositoryImpl) SaveOperationGroup(tx *gorm.DB, req masteroperationpayloads.OperationGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationGroup{
		IsActive:                  req.IsActive,
		OperationGroupId:          req.OperationGroupId,
		OperationGroupCode:        req.OperationGroupCode,
		OperationGroupDescription: req.OperationGroupDescription,
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

	if len(req.OperationGroupCode) > 2 {
		// errMessage := "Operation Group Code max 2 characters"

		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,

			Err: errors.New(utils.BadRequestError),
		}
	}

	return true, nil
}

func (*OperationGroupRepositoryImpl) ChangeStatusOperationGroup(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteroperationentities.OperationGroup

	result := tx.Model(&entities).
		Where("operation_group_id = ?", Id).
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
