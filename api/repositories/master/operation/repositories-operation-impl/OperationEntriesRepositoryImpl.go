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

	"gorm.io/gorm"
)

type OperationEntriesRepositoryImpl struct {
}

func StartOperationEntriesRepositoryImpl() masteroperationrepository.OperationEntriesRepository {
	return &OperationEntriesRepositoryImpl{}
}

func (r *OperationEntriesRepositoryImpl) GetAllOperationEntries(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteroperationentities.OperationEntries{}
	var responses []masteroperationpayloads.OperationEntriesResponse

	// define table struct
	tableStruct := masteroperationpayloads.OperationEntriesResponse{}

	//join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []masteroperationentities.OperationEntries{}
		return pages, nil
	}

	pages.Rows = responses

	return pages, nil
}

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesName(tx *gorm.DB, request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse) {
	tableStruct := masteroperationpayloads.OperationEntriesResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", request.OperationGroupId).
		Where("mtr_operation_section.operation_section_id = ?", request.OperationSectionId).
		Where("mtr_operation_key.operation_key_id = ?", request.OperationKeyId).
		Where("mtr_operation_entries.operation_entries_code = ?", request.OperationEntriesCode)

	rows, err := WhereQuery.First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return tableStruct, nil
}

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationEntriesResponse, *exceptions.BaseErrorResponse) {
	response := masteroperationpayloads.OperationEntriesResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, response)

	whereQuery := joinTable.Where("operation_entries_id = ?", Id)

	rows, err := whereQuery.First(&response).Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

// func (r *OperationEntriesRepositoryImpl) GetOperationEntriesKeyCodeByGroupId

func (r *OperationEntriesRepositoryImpl) SaveOperationEntries(tx *gorm.DB, request masteroperationpayloads.OperationEntriesResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationEntries{
		IsActive:                    request.IsActive,
		OperationEntriesId:          request.OperationEntriesId,
		OperationEntriesCode:        request.OperationEntriesCode,
		OperationGroupId:            request.OperationGroupId,
		OperationSectionId:          request.OperationSectionId,
		OperationKeyId:              request.OperationKeyId,
		OperationEntriesDescription: request.OperationEntriesDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(request.OperationEntriesCode) > 6 {
		// errMessage := "Operation Entries Code max 6 characters"

		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,

			Err: errors.New(utils.BadRequestError),
		}
	}

	return true, nil
}

func (r *OperationEntriesRepositoryImpl) ChangeStatusOperationEntries(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationEntries
	// var response masteroperationpayloads.OperationEntriesResponse
	result := tx.Model(&entities).
		Where("operation_entries_id = ?", Id).
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
