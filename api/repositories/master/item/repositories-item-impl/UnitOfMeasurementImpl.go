package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepositoryImpl struct {
	myDB *gorm.DB
}

func StartUnitOfMeasurementRepositoryImpl(db *gorm.DB) masteritemrepository.UnitOfMeasurementRepository {
	return &UnitOfMeasurementRepositoryImpl{myDB: db}
}

func (r *UnitOfMeasurementRepositoryImpl) WithTrx(trxHandle *gorm.DB) masteritemrepository.UnitOfMeasurementRepository {
	if trxHandle == nil {
		log.Println("Transaction Database Not Found!")
		return r
	}
	r.myDB = trxHandle
	return r
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	entities := masteritementities.Uom{}
	var responses []masteritempayloads.UomResponse
	// define table struct
	tableStruct := masteritempayloads.UomResponse{}
	//define join table
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

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, error) {
	var UnitOfMeasurements []masteritementities.Uom
	response := []masteritempayloads.UomResponse{}

	err := r.myDB.Model(&UnitOfMeasurements).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementById(Id int) (masteritempayloads.UomResponse, error) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomResponse, error) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomResponse{}

	rows, err := r.myDB.Model(&entities).
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

func (r *UnitOfMeasurementRepositoryImpl) SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, error) {
	entities := masteritementities.Uom{
		IsActive:       req.IsActive,
		UomId:          req.UomId,
		UomTypeId:      req.UomTypeId,
		UomCode:        req.UomCode,
		UomDescription: req.UomDescription,
	}

	err := r.myDB.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UnitOfMeasurementRepositoryImpl) ChangeStatusUnitOfMeasurement(Id int) (bool, error) {
	var entities masteritementities.Uom

	result := r.myDB.Model(&entities).
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

	result = r.myDB.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}
