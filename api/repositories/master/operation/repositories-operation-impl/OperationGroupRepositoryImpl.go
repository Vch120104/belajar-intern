package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type OperationGroupRepositoryImpl struct {
}

func StartOperationGroupRepositoryImpl() masteroperationrepository.OperationGroupRepository {
	return &OperationGroupRepositoryImpl{}
}

func (r *OperationGroupRepositoryImpl) GetAllOperationGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masteroperationentities.OperationGroup{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if len(entities) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *OperationGroupRepositoryImpl) GetAllOperationGroupIsActive(tx *gorm.DB) ([]masteroperationpayloads.OperationGroupResponse, error) {
	var OperationGroups []masteroperationentities.OperationGroup
	response := []masteroperationpayloads.OperationGroupResponse{}

	err := tx.Model(&OperationGroups).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *OperationGroupRepositoryImpl) GetOperationGroupById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationGroupResponse, error) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationGroupRepositoryImpl) GetOperationGroupByCode(tx *gorm.DB, Code string) (masteroperationpayloads.OperationGroupResponse, error) {
	entities := masteroperationentities.OperationGroup{}
	response := masteroperationpayloads.OperationGroupResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationGroup{
			OperationGroupCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationGroupRepositoryImpl) SaveOperationGroup(tx *gorm.DB, req masteroperationpayloads.OperationGroupResponse) (bool, error) {
	entities := masteroperationentities.OperationGroup{
		IsActive:                  req.IsActive,
		OperationGroupId:          req.OperationGroupId,
		OperationGroupCode:        req.OperationGroupCode,
		OperationGroupDescription: req.OperationGroupDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OperationGroupRepositoryImpl) ChangeStatusOperationGroup(tx *gorm.DB, Id int) (bool, error) {
	var entities masteroperationentities.OperationGroup

	result := tx.Model(&entities).
		Where("operation_group_id = ?", Id).
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
