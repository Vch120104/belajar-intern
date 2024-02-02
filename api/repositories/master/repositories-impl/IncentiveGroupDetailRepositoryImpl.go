package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"

	"log"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepositoryImpl struct {
}

func StartIncentiveGroupDetailRepositoryImpl() masterrepository.IncentiveGroupDetailRepository {
	return &IncentiveGroupDetailRepositoryImpl{}
}

func (r *IncentiveGroupDetailRepositoryImpl) GetAllIncentiveGroupDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.IncentiveGroupDetail{}
	//define base model
	baseModelQuery := r.myDB.Model(&entities)
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

func (r *IncentiveGroupDetailRepositoryImpl) GetAlIncentiveGroupDetailIsActive() ([]masterpayloads.IncentiveGroupDetailResponse, error) {
	var IncentiveGroupsDetail []masterentities.IncentiveGroupDetail
	response := []masterpayloads.IncentiveGroupDetailResponse{}

	err := r.myDB.Model(&IncentiveGroupsDetail).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) GetIncentiveGroupDetailById(Id int) (masterpayloads.IncentiveGroupDetailResponse, error) {
	entities := masterentities.IncentiveGroupDetail{}
	response := masterpayloads.IncentiveGroupDetailResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *IncentiveGroupDetailRepositoryImpl) GetIncentiveGroupDetailByCode(Code string) (masterpayloads.IncentiveGroupDetailResponse, error) {
	entities := masterentities.IncentiveGroupDetail{}
	response := masterpayloads.IncentiveGroupDetailResponse{}

	rows, err := r.myDB.Model(&entities).
		Where(masterentities.IncentiveGroupDetail{
			IncentiveGroupCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailResponse) (bool, error) {
	entities := masterentities.IncentiveGroupDetail{
		IncentiveGroupDetailId : req.IncentiveGroupDetailId,
		IncentiveGroupId : req.IncentiveGroupId,
		IncentiveGroupCode : req.IncentiveGroupCode,  
		IncentiveLevel : req.IncentiveLevel,
		TargetAmount : req.TargetAmount,
		TargetPercent : req.TargetPercent,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *IncentiveGroupDetailRepositoryImpl) ChangeStatusIncentiveGroupDetail(Id int) (bool, error) {
	var entities masterentities.IncentiveGroupDetail

	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
