package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OperationKeyRepositoryImpl struct {
}

func StartOperationKeyRepositoryImpl() masteroperationrepository.OperationKeyRepository {
	return &OperationKeyRepositoryImpl{}
}

func (r *OperationKeyRepositoryImpl) GetAllOperationKeyList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteroperationentities.OperationKey{}
	var responses []masteroperationpayloads.OperationkeyListResponse

	// define table struct
	tableStruct := masteroperationpayloads.OperationkeyListResponse{}

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

func (r *OperationKeyRepositoryImpl) GetOperationKeyById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationkeyListResponse, error) {
	response := masteroperationpayloads.OperationkeyListResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, response)

	whereQuery := joinTable.Where("operation_key_id = ?", Id)

	rows, err := whereQuery.First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

// func (r *OperationKeyRepositoryImpl) GetOperationKeyCode(request masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyCodeResponse, error) {
// 	entities := masteroperationentities.OperationKey{}
// 	response := masteroperationpayloads.OperationKeyCodeResponse{}

// 	rows, err := r.myDB.Model(&entities).
// 		Where(masteroperationpayloads.OperationKeyCodeResponse{
// 			OperationGroupId:   int(request.OperationGroupId),
// 			OperationSectionId: int(request.OperationSectionId),
// 		}).
// 		First(&response).
// 		Rows()

// 	if err != nil {
// 		return response, err
// 	}

// 	defer rows.Close()

// 	return response, nil
// }

func (r *OperationKeyRepositoryImpl) GetOperationKeyName(tx *gorm.DB, request masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error) {
	tableStruct := masteroperationpayloads.OperationKeyNameResponse{}
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	tx.Logger = newLogger

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	WhereQuery := joinTable.
		Where("mtr_operation_group.operation_group_id = ?", request.OperationGroupId).
		Where("mtr_operation_section.operation_section_id = ?", request.OperationSectionId).
		Where("mtr_operation_key.operation_key_code = ?", request.OperationKeyCode)

	rows, err := WhereQuery.First(&tableStruct).Rows()

	if err != nil {
		return tableStruct, err
	}

	defer rows.Close()

	return tableStruct, nil
}

func (r *OperationKeyRepositoryImpl) SaveOperationKey(tx *gorm.DB, request masteroperationpayloads.OperationKeyResponse) (bool, error) {
	entities := masteroperationentities.OperationKey{
		IsActive:                request.IsActive,
		OperationKeyId:          request.OperationKeyId,
		OperationKeyCode:        request.OperationKeyCode,
		OperationGroupId:        request.OperationGroupId,
		OperationSectionId:      request.OperationSectionId,
		OperationKeyDescription: request.OperationKeyDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationKeyRepositoryImpl) ChangeStatusOperationKey(tx *gorm.DB, Id int) (bool, error) {
	var entities masteroperationentities.OperationKey

	result := tx.Model(&entities).
		Where("operation_key_id = ?", Id).
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
