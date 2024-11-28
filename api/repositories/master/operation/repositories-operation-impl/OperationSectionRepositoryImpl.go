package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"math"
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
	var entities []masteroperationentities.OperationSection
	var responses []masteroperationpayloads.OperationSectionListResponse

	tableStruct := masteroperationpayloads.OperationSectionListResponse{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.Limit)))

	if err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).
		Find(&entities).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = responses
		return pages, nil
	}

	var operationGroupIds []int
	for _, entity := range entities {
		operationGroupIds = append(operationGroupIds, entity.OperationGroupId)
	}

	var operationGroups []masteroperationentities.OperationGroup
	if err := tx.Where("operation_group_id IN ?", operationGroupIds).Find(&operationGroups).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	operationGroupMap := make(map[int]masteroperationentities.OperationGroup)
	for _, group := range operationGroups {
		operationGroupMap[group.OperationGroupId] = group
	}

	for _, entity := range entities {
		response := masteroperationpayloads.OperationSectionListResponse{
			IsActive:                    entity.IsActive,
			OperationSectionId:          entity.OperationSectionId,
			OperationSectionCode:        entity.OperationSectionCode,
			OperationSectionDescription: entity.OperationSectionDescription,
			OperationGroupId:            entity.OperationGroupId,
		}

		if operationGroup, exists := operationGroupMap[entity.OperationGroupId]; exists {
			response.OperationGroupCode = operationGroup.OperationGroupCode
			response.OperationGroupDescription = operationGroup.OperationGroupDescription
		}

		responses = append(responses, response)
	}

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

func (r *OperationSectionRepositoryImpl) GetOperationSectionDropDown(tx *gorm.DB, operationGroupId int) ([]masteroperationpayloads.OperationSectionDropDown, *exceptions.BaseErrorResponse) {

	var OperationSectionDropDown []masteroperationpayloads.OperationSectionDropDown

	err := tx.Model(&masteroperationentities.OperationSection{}).
		Select("operation_section_id", "CONCAT(operation_section_code, ' - ', operation_section_description) as operation_section_code").
		Where(masteroperationentities.OperationGroup{OperationGroupId: operationGroupId}).
		Scan(&OperationSectionDropDown)
	if err.Error != nil {
		return OperationSectionDropDown, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error,
		}
	}
	return OperationSectionDropDown, nil
}
