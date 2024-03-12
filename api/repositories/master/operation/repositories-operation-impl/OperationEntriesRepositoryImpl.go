package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationEntriesRepositoryImpl struct {
}

func StartOperationEntriesRepositoryImpl() masteroperationrepository.OperationEntriesRepository {
	return &OperationEntriesRepositoryImpl{}
}

func (r *OperationEntriesRepositoryImpl) GetAllOperationEntries(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteroperationentities.OperationEntries{}
	var responses []masteroperationpayloads.OperationEntriesResponse

	// define table struct
	tableStruct := masteroperationpayloads.OperationEntriesResponse{}

	//join table
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

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesName(tx *gorm.DB, request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error) {
	tableStruct := masteroperationpayloads.OperationEntriesResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", request.OperationGroupId).
		Where("mtr_operation_section.operation_section_id = ?", request.OperationSectionId).
		Where("mtr_operation_key.operation_key_id = ?", request.OperationKeyId).
		Where("mtr_operation_entries.operation_entries_code = ?", request.OperationEntriesCode)

	rows, err := WhereQuery.First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, err
	}

	defer rows.Close()

	return tableStruct, nil
}

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationEntriesResponse, error) {
	response := masteroperationpayloads.OperationEntriesResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, response)

	whereQuery := joinTable.Where("operation_entries_id = ?", Id)

	rows, err := whereQuery.First(&response).Rows()


	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

// func (r *OperationEntriesRepositoryImpl) GetOperationEntriesKeyCodeByGroupId

func (r *OperationEntriesRepositoryImpl) SaveOperationEntries(tx *gorm.DB, request masteroperationpayloads.OperationEntriesResponse) (bool, error) {
	entities := masteroperationentities.OperationEntries{
		IsActive:             request.IsActive,
		OperationEntriesId:   request.OperationEntriesId,
		OperationEntriesCode: request.OperationEntriesCode,
		OperationGroupId:     request.OperationGroupId,
		OperationSectionId:   request.OperationSectionId,
		OperationKeyId:       request.OperationKeyId,
		OperationEntriesDesc: request.OperationEntriesDesc,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationEntriesRepositoryImpl) ChangeStatusOperationEntries(tx *gorm.DB, Id int) (bool, error) {
	var entities masteroperationentities.OperationEntries
	result := tx.Model(&entities).
		Where("operation_entries_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

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
