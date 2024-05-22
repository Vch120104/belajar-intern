package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepositoryImpl struct {
}

func StartUnitOfMeasurementRepositoryImpl() masteritemrepository.UnitOfMeasurementRepository {
	return &UnitOfMeasurementRepositoryImpl{}
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurementIsActive(tx *gorm.DB) ([]masteritempayloads.UomResponse, *exceptions.BaseErrorResponse) {
	var UnitOfMeasurements []masteritementities.Uom
	response := []masteritempayloads.UomResponse{}

	err := tx.Model(&UnitOfMeasurements).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementById(tx *gorm.DB, Id int) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomIdCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementByCode(tx *gorm.DB, Code string) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomIdCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) SaveUnitOfMeasurement(tx *gorm.DB, req masteritempayloads.UomResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.Uom{
		IsActive:       req.IsActive,
		UomId:          req.UomId,
		UomTypeId:      req.UomTypeId,
		UomCode:        req.UomCode,
		UomDescription: req.UomDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *UnitOfMeasurementRepositoryImpl) ChangeStatusUnitOfMeasurement(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.Uom

	result := tx.Model(&entities).
		Where("uom_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
