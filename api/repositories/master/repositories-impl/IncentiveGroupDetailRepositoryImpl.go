package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepositoryImpl struct {
}

func StartIncentiveGroupDetailRepositoryImpl() masterrepository.IncentiveGroupDetailRepository {
	return &IncentiveGroupDetailRepositoryImpl{}
}

func (r *IncentiveGroupDetailRepositoryImpl) GetAllIncentiveGroupDetail(tx *gorm.DB, headerId int, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.IncentiveGroupDetail{}
	response := []masterpayloads.IncentiveGroupDetailResponse{}
	//define base model
	baseModelQuery := tx.Model(&entities).Where("mtr_incentive_group_detail.incentive_group_id = ?", headerId)

	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&response).Rows()

	if len(response) == 0 {
		return pages, gorm.ErrRecordNotFound
	}

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = response

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

func (r *IncentiveGroupDetailRepositoryImpl) SaveIncentiveGroupDetail(tx *gorm.DB, req masterpayloads.IncentiveGroupDetailRequest) (bool, error) {
	entities := masterentities.IncentiveGroupDetail{
		IsActive:               req.IsActive,
		IncentiveGroupDetailId: req.IncentiveGroupDetailId,
		IncentiveGroupId:       req.IncentiveGroupId,
		IncentiveLevel:         req.IncentiveLevel,
		TargetAmount:           req.TargetAmount,
		TargetPercent:          req.TargetPercent,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
