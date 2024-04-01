package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
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

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
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
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = responses

	return pages, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurementIsActive(tx *gorm.DB) ([]masteritempayloads.UomResponse, *exceptionsss_test.BaseErrorResponse) {
	var UnitOfMeasurements []masteritementities.Uom
	response := []masteritempayloads.UomResponse{}

	err := tx.Model(&UnitOfMeasurements).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementById(tx *gorm.DB,Id int) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomIdCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementByCode(tx *gorm.DB,Code string) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.Uom{}
	response := masteritempayloads.UomIdCodeResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.Uom{
			UomCode: Code,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *UnitOfMeasurementRepositoryImpl) SaveUnitOfMeasurement(tx *gorm.DB,req masteritempayloads.UomResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *UnitOfMeasurementRepositoryImpl) ChangeStatusUnitOfMeasurement(tx *gorm.DB,Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.Uom

	result := tx.Model(&entities).
		Where("uom_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
