package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type OperationEntriesRepositoryImpl struct {
	myDB *gorm.DB
}

func StartOperationEntriesRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationEntriesRepository {
	return &OperationEntriesRepositoryImpl{myDB: db}
}

func (r *OperationEntriesRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationEntriesRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesName(request masteroperationpayloads.OperationEntriesRequest) (masteroperationpayloads.OperationEntriesResponse, error) {
	tableStruct := masteroperationpayloads.OperationEntriesResponse{}

	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

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

func (r *OperationEntriesRepositoryImpl) GetOperationEntriesById(Id int32) (masteroperationpayloads.OperationEntriesResponse, error) {
	entities := masteroperationentities.OperationEntries{}
	response := masteroperationpayloads.OperationEntriesResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationEntries{
			OperationEntriesId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

// func (r *OperationEntriesRepositoryImpl) GetOperationEntriesKeyCodeByGroupId

func (r *OperationEntriesRepositoryImpl) SaveOperationEntries(request masteroperationpayloads.OperationEntriesResponse) (bool, error) {
	entities := masteroperationentities.OperationEntries{
		IsActive:             request.IsActive,
		OperationEntriesId:   request.OperationEntriesId,
		OperationEntriesCode: request.OperationEntriesCode,
		OperationGroupId:     request.OperationGroupId,
		OperationSectionId:   request.OperationSectionId,
		OperationKeyId:       request.OperationKeyId,
		OperationEntriesDesc: request.OperationEntriesDesc,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationEntriesRepositoryImpl) ChangeStatusOperationEntries(Id int) (bool, error) {
	var entities masteroperationentities.OperationEntries
	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
