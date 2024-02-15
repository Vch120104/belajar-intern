package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepositoryImpl struct {
}

func StartUnitOfMeasurementRepositoryImpl() masteritemrepository.UnitOfMeasurementRepository {
	return &UnitOfMeasurementRepositoryImpl{}
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteritementities.Uom{}
	var responses []masteritempayloads.UomResponse
	// define table struct
	tableStruct := masteritempayloads.UomResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
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

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurementIsActive(tx *gorm.DB) ([]masteritempayloads.UomResponse, error) {
	var UnitOfMeasurements []masteritementities.Uom
	response := []masteritempayloads.UomResponse{}

	err := tx.Model(&UnitOfMeasurements).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementById(tx *gorm.DB,Id int) (masteritempayloads.UomResponse, error) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementByCode(tx *gorm.DB,Code string) (masteritempayloads.UomResponse, error) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, err
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) SaveUnitOfMeasurement(tx *gorm.DB,req masteritempayloads.UomResponse) (bool, error) {
	entities := masteritementities.Uom{
		IsActive:       req.IsActive,
		UomId:          req.UomId,
		UomTypeId:      req.UomTypeId,
		UomCode:        req.UomCode,
		UomDescription: req.UomDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UnitOfMeasurementRepositoryImpl) ChangeStatusUnitOfMeasurement(tx *gorm.DB,Id int) (bool, error) {
	var entities masteritementities.Uom

	result := tx.Model(&entities).
		Where("uom_id = ?", Id).
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
