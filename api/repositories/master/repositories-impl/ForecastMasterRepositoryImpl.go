package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"

	"gorm.io/gorm"
)

type ForecastMasterRepositoryImpl struct {
}

func StartForecastMasterRepositoryImpl() masterrepository.ForecastMasterRepository {
	return &ForecastMasterRepositoryImpl{}
}

func (r *ForecastMasterRepositoryImpl) GetForecastMasterById(tx *gorm.DB, Id int) (masterpayloads.ForecastMasterResponse, error) {
	entities := masterentities.ForecastMaster{}
	response := masterpayloads.ForecastMasterResponse{}

	rows, err := tx.Model(&entities).
		Where(masterentities.ForecastMaster{
			ForecastMasterId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *ForecastMasterRepositoryImpl) SaveForecastMaster(tx *gorm.DB, req masterpayloads.ForecastMasterResponse) (bool, error) {
	entities := masterentities.ForecastMaster{
		IsActive:                   req.IsActive,
		ForecastMasterId:           req.ForecastMasterId,
		SupplierId:                 req.SupplierId,
		MovingCodeId:               req.MovingCodeId,
		OrderTypeId:                req.OrderTypeId,
		ForecastMasterLeadTime:     req.ForecastMasterLeadTime,
		ForecastMasterSafetyFactor: req.ForecastMasterSafetyFactor,
		ForecastMasterOrderCycle:   req.ForecastMasterOrderCycle,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ForecastMasterRepositoryImpl) ChangeStatusForecastMaster(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.ForecastMaster

	result := tx.Model(&entities).
		Where("forecast_master_id = ?", Id).
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
