package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupRepositoryImpl struct {
}

func StartIncentiveGroupRepositoryImpl() masterrepository.IncentiveGroupRepository {
	return &IncentiveGroupRepositoryImpl{}
}

func (r *IncentiveGroupRepositoryImpl) GetAllIncentiveGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.IncentiveGroup{}
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

func (r *IncentiveGroupRepositoryImpl) GetAllIncentiveGroupIsActive(tx *gorm.DB) ([]masterpayloads.IncentiveGroupResponse, error) {
	var IncentiveGroups []masterentities.IncentiveGroup
	response := []masterpayloads.IncentiveGroupResponse{}

	err := tx.Model(&IncentiveGroups).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (*IncentiveGroupRepositoryImpl) GetIncentiveGroupById(tx *gorm.DB, Id int) (masterpayloads.IncentiveGroupResponse, error) {
	entities := masterentities.IncentiveGroup{}
	response := masterpayloads.IncentiveGroupResponse{}

	rows, err := tx.Model(&entities).Where(masterentities.IncentiveGroup{
		IncentiveGroupId: Id,
	}).First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()
	return response, nil
}

func (r *IncentiveGroupRepositoryImpl) SaveIncentiveGroup(tx *gorm.DB, req masterpayloads.IncentiveGroupResponse) (bool, error) {
	entities := masterentities.IncentiveGroup{
		IsActive:           req.IsActive,
		IncentiveGroupId:   req.IncentiveGroupId,
		IncentiveGroupCode: req.IncentiveGroupCode,
		IncentiveGroupName: req.IncentiveGroupName,
		EffectiveDate:      req.EffectiveDate,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *IncentiveGroupRepositoryImpl) ChangeStatusIncentiveGroup(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.IncentiveGroup

	result := tx.Model(&entities).
		Where("incentive_group_id = ?", Id).
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
