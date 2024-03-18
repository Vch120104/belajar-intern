	package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationSectionRepositoryImpl struct {
}

func StartOperationSectionRepositoryImpl() masteroperationrepository.OperationSectionRepository {
	return &OperationSectionRepositoryImpl{}
}

func (r *OperationSectionRepositoryImpl) GetAllOperationSectionList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
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
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil

}

func (r *OperationSectionRepositoryImpl) GetOperationSectionName(tx *gorm.DB, GroupId int, SectionCode string) (masteroperationpayloads.OperationSectionNameResponse, error) {
	tableStruct := masteroperationpayloads.OperationSectionNameResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	row, err := joinTable.Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.operation_section_code = ?", SectionCode).
		First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, err
	}

	defer row.Close()

	return tableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetSectionCodeByGroupId(tx *gorm.DB, GroupId int) ([]masteroperationpayloads.OperationSectionCodeResponse, error) {
	tableStruct := masteroperationpayloads.OperationSectionCodeResponse{}
	var sliceTableStruct []masteroperationpayloads.OperationSectionCodeResponse

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", GroupId).
		Where("mtr_operation_section.is_active = 1")

	rows, err := WhereQuery.Scan(&sliceTableStruct).Rows()

	if len(sliceTableStruct) == 0 {
		return sliceTableStruct, gorm.ErrRecordNotFound
	}

	if err != nil {
		return sliceTableStruct, err
	}
	defer rows.Close()

	return sliceTableStruct, nil
}

func (r *OperationSectionRepositoryImpl) GetOperationSectionById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationSectionListResponse, error) {
	response := masteroperationpayloads.OperationSectionListResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, response)

	whereQuery := joinTable.Where("operation_section_id = ?", Id)

	rows, err := whereQuery.First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationSectionRepositoryImpl) SaveOperationSection(tx *gorm.DB, request masteroperationpayloads.OperationSectionRequest) (bool, error) {
	entities := masteroperationentities.OperationSection{
		IsActive:                    request.IsActive,
		OperationSectionId:          request.OperationSectionId,
		OperationSectionCode:        request.OperationSectionCode,
		OperationGroupId:            request.OperationGroupId,
		OperationSectionDescription: request.OperationSectionDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationSectionRepositoryImpl) ChangeStatusOperationSection(tx *gorm.DB, Id int) (bool, error) {
	var entities masteroperationentities.OperationSection

	result := tx.Model(&entities).
		Where("operation_section_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	// Toggle the IsActive value
	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
