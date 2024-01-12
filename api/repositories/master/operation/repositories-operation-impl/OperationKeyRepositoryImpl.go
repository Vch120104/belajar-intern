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
	myDB *gorm.DB
}

func StartOperationKeyRepositoryImpl(db *gorm.DB) masteroperationrepository.OperationKeyRepository {
	return &OperationKeyRepositoryImpl{myDB: db}
}

func (r *OperationKeyRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteroperationrepository.OperationKeyRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *OperationKeyRepositoryImpl) GetAllOperationKeyList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteroperationentities.OperationKey{}
	var responses []masteroperationpayloads.OperationkeyListResponse

	// define table struct
	tableStruct := masteroperationpayloads.OperationkeyListResponse{}

	//join table
	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	//apply pagination and execute
	rows, err := joinTable.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Rows()

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *OperationKeyRepositoryImpl) GetOperationKeyById(Id int) (masteroperationpayloads.OperationKeyResponse, error) {
	entities := masteroperationentities.OperationKey{}
	response := masteroperationpayloads.OperationKeyResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masteroperationentities.OperationKey{
			OperationKeyId: Id,
		}).
		First(&response).
		Rows()

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

func (r *OperationKeyRepositoryImpl) GetOperationKeyName(request masteroperationpayloads.OperationKeyRequest) (masteroperationpayloads.OperationKeyNameResponse, error) {
	tableStruct := masteroperationpayloads.OperationKeyNameResponse{}
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	r.myDB.Logger = newLogger

	joinTable := utils.CreateJoinSelectStatement(r.myDB, tableStruct)

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

func (r *OperationKeyRepositoryImpl) SaveOperationKey(request masteroperationpayloads.OperationKeyResponse) (bool, error) {
	entities := masteroperationentities.OperationKey{
		IsActive:                request.IsActive,
		OperationKeyId:          request.OperationKeyId,
		OperationKeyCode:        request.OperationKeyCode,
		OperationGroupId:        request.OperationGroupId,
		OperationSectionId:      request.OperationSectionId,
		OperationKeyDescription: request.OperationKeyDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationKeyRepositoryImpl) ChangeStatusOperationKey(Id int) (bool, error) {
	var entities masteroperationentities.OperationKey

	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
