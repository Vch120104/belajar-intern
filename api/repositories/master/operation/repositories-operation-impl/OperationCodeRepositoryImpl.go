package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationCodeRepositoryImpl struct {
}

func StartOperationCodeRepositoryImpl() masteroperationrepository.OperationCodeRepository {
	return &OperationCodeRepositoryImpl{}
}

func (r *OperationCodeRepositoryImpl) GetAllOperationCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteroperationentities.OperationCode{}
	var payloads []masteroperationpayloads.OperationCodeGetAll
	baseModelQuery := tx.Model(&entities)
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&payloads).Rows()
	if len(payloads) == 0 {
		return pages, gorm.ErrRecordNotFound
	}
	if err != nil {
		return pages, err
	}
	defer rows.Close()

	pages.Rows = payloads

	return pages, nil
}

func (r *OperationCodeRepositoryImpl) GetOperationCodeById(tx *gorm.DB, id int) (masteroperationpayloads.OperationCodeResponse, error) {
	entities := masteroperationentities.OperationCode{}
	response := masteroperationpayloads.OperationCodeResponse{}
	rows, err := tx.Model(&entities).Where(masteroperationentities.OperationCode{OperationId: id}).First(&response).Rows()
	if err != nil {
		return response, err
	}
	defer rows.Close()
	return response, nil
}

func (r *OperationCodeRepositoryImpl) SaveOperationCode(tx *gorm.DB, req masteroperationpayloads.OperationCodeSave) (bool, error) {
	entities := masteroperationentities.OperationCode{
		IsActive:                req.IsActive,
		OperationCode:           req.OperationCode,
		OperationName:           req.OperationName,
		OperationUsingIncentive: req.OperationUsingIncentive,
		OperationUsingActual:    req.OperationUsingActual,
	}
	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *OperationCodeRepositoryImpl) ChangeStatusItemSubstitute(tx *gorm.DB, id int) (bool, error) {
	var entities masteroperationentities.OperationCode

	result := tx.Model(&entities).
		Where("operation_id = ?", id).
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
