package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"

	"gorm.io/gorm"
)

type MovingCodeRepositoryImpl struct {
}

func StartMovingCodeRepositoryImpl() masterrepository.MovingCodeRepository {
	return &MovingCodeRepositoryImpl{}
}

func (*MovingCodeRepositoryImpl) GetAllMovingCode(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := []masterentities.MovingCode{}
	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply pagination and execute
	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&entities).Rows()

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

func (*MovingCodeRepositoryImpl) GetMovingCodeById(tx *gorm.DB, Id int) (masterpayloads.MovingCodeResponse, error) {
	entities := masterentities.MovingCode{}
	response := masterpayloads.MovingCodeResponse{}

	rows, err := tx.Model(&entities).Where(masterentities.MovingCode{
		MovingCodeId: Id,
	}).First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()
	return response, nil
}

func (*MovingCodeRepositoryImpl) GetMovingCodeByPriority(tx *gorm.DB, Priority float64) (masterpayloads.MovingCodeResponse, error) {
	entities := masterentities.MovingCode{}
	response := masterpayloads.MovingCodeResponse{}

	rows, err := tx.Model(&entities).Where(masterentities.MovingCode{
		Priority: Priority,
	}).First(&response).Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()
	return response, nil
}

func (r *MovingCodeRepositoryImpl) SaveMovingCode(tx *gorm.DB, req masterpayloads.MovingCodeRequest) (bool, error) {
	entities := masterentities.MovingCode{
		IsActive:              req.IsActive,
		CompanyId:             req.CompanyId,
		MovingCodeDescription: req.MovingCodeDescription,
		DemandExistMonthFrom:  req.DemandExistMonthFrom,
		DemandExistMonthTo:    req.DemandExistMonthTo,
		AgingMonthFrom:        req.AgingMonthFrom,
		AgingMonthTo:          req.AgingMonthTo,
		LastMovingMonthFrom:   req.LastMovingMonthFrom,
		LastMovingMonthTo:     req.LastMovingMonthTo,
		MinimumQuantityDemand: req.MinimumQuantityDemand,
		Remark:                req.Remark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *MovingCodeRepositoryImpl) IncreasePriorityMovingCode(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.MovingCode

	result := tx.Model(&entities).
		Where("moving_code_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	entities.Priority += 1

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *MovingCodeRepositoryImpl) DecreasePriorityMovingCode(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.MovingCode

	result := tx.Model(&entities).
		Where("moving_code_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.Priority > 1 {
		entities.Priority -= 1
	} else {
		entities.Priority -= 0
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *MovingCodeRepositoryImpl) ChangeStatusMovingCode(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.MovingCode

	result := tx.Model(&entities).
		Where("moving_code_id = ?", Id).
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
