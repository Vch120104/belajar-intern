package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"net/http"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationSectionRepositoryImpl struct {
}

func StartOperationSectionRepositoryImpl() masteroperationrepository.OperationSectionRepository {
	return &OperationSectionRepositoryImpl{}
}

func (r *OperationSectionRepositoryImpl) GetAllOperationSectionList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationSection{}
	var responses []masteroperationpayloads.OperationSectionListResponse
	// define table struct
	tableStruct := masteroperationpayloads.OperationSectionListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	//apply pagination and execute
	rows, err := joinTable.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        err,
		}
	}

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil

}

func (r *OperationSectionRepositoryImpl) GetOperationSectionName(tx *gorm.DB, GroupId int, SectionCode string) (masteroperationpayloads.OperationSectionNameResponse, *exceptions.BaseErrorResponse) {
	tableStruct := masteroperationpayloads.OperationSectionNameResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	row, err := joinTable.Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.operation_section_code = ?", SectionCode).
		First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer row.Close()

	return tableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetSectionCodeByGroupId(tx *gorm.DB, GroupId int) ([]masteroperationpayloads.OperationSectionCodeResponse, *exceptions.BaseErrorResponse) {
	tableStruct := masteroperationpayloads.OperationSectionCodeResponse{}
	var sliceTableStruct []masteroperationpayloads.OperationSectionCodeResponse

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.is_active = 1")

	rows, err := WhereQuery.Scan(&sliceTableStruct).Rows()

	if len(sliceTableStruct) == 0 {
		return sliceTableStruct, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        err,
		}
	}

	if err != nil {
		return sliceTableStruct, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return sliceTableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetOperationSectionById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationSectionListResponse, *exceptions.BaseErrorResponse) {
	response := masteroperationpayloads.OperationSectionListResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, response)

	whereQuery := joinTable.Where("operation_section_id = ?", Id)

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

func (r *OperationSectionRepositoryImpl) SaveOperationSection(tx *gorm.DB, request masteroperationpayloads.OperationSectionRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationSection{
		IsActive:                    request.IsActive,
		OperationSectionId:          request.OperationSectionId,
		OperationSectionCode:        request.OperationSectionCode,
		OperationGroupId:            request.OperationGroupId,
		OperationSectionDescription: request.OperationSectionDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationSectionRepositoryImpl) ChangeStatusOperationSection(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationSection

	result := tx.Model(&entities).
		Where("operation_section_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	// Toggle the IsActive value
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
