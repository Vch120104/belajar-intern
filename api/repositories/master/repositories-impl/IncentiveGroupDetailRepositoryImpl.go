package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepositoryImpl struct {
}

func StartIncentiveGroupDetailRepositoryImpl() masterrepository.IncentiveGroupDetailRepository {
	return &IncentiveGroupDetailRepositoryImpl{}
}

func (r *IncentiveGroupDetailRepositoryImpl) GetAllIncentiveGroupDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.IncentiveGroupDetail{}
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

func (r *IncentiveGroupDetailRepositoryImpl) GetIncentiveGroupDetailById(tx *gorm.DB, Id int) (masterpayloads.IncentiveGroupDetailResponse, error) {
	entities := masterentities.IncentiveGroupDetail{}
	response := masterpayloads.IncentiveGroupDetailResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.IncentiveGroupDetail{
			IncentiveGroupDetailId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) SaveIncentiveGroupDetail(tx *gorm.DB, IncentiveGroupId int, req masterpayloads.IncentiveGroupDetailResponse) (bool, error) {
	entities := masterentities.IncentiveGroupDetail{
		IncentiveGroupDetailId : req.IncentiveGroupDetailId,
		IncentiveGroupId : IncentiveGroupId,
		IncentiveLevel : req.IncentiveLevel,
		TargetAmount : req.TargetAmount,
		TargetPercent : req.TargetPercent,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) ChangeStatusIncentiveGroupDetail(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.IncentiveGroupDetail

	result := tx.Model(&entities).
		Where("incentive_group_detail_id = ?", Id).
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
