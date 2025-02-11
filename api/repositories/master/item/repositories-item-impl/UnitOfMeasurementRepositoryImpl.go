package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepositoryImpl struct {
}

func (r *UnitOfMeasurementRepositoryImpl) GetUnitOfMeasurementItem(tx *gorm.DB, Payload masteritempayloads.UomItemRequest) (masteritempayloads.UomItemResponses, *exceptions.BaseErrorResponse) {
	response := masteritempayloads.UomItemResponses{}

	// rows, err := tx.Model(&entities).
	// 	Where(masteritementities.UomItem{ItemId: Payload.ItemId, UomSourceTypeCode: Payload.SourceType}).
	// 	First(&response).
	// 	Rows()

	// if err != nil {
	// 	return response, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        err,
	// 	}
	// } else {
	// 	defer func(rows *sql.Rows) {
	// 		err := rows.Close()
	// 		if err != nil {

	// 		}
	// 	}(rows)
	// }

	// refactored empty branch
	entities := masteritementities.UomItem{}
	err := tx.Model(&entities).
		Where(masteritementities.UomItem{ItemId: Payload.ItemId, UomSourceTypeCode: Payload.SourceType}).
		First(&response).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item not found",
				Err:        fmt.Errorf("item not found"),
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error, please check entity",
			Err:        err,
		}
	}

	return response, nil
}

func StartUnitOfMeasurementRepositoryImpl() masteritemrepository.UnitOfMeasurementRepository {
	return &UnitOfMeasurementRepositoryImpl{}
}

func (r *UnitOfMeasurementRepositoryImpl) GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.UomResponse
	tableStruct := masteritempayloads.UomResponse{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.Limit)))

	if err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).
		Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []masteritempayloads.UomResponse{}
		return pages, nil
	}

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
func (r *UnitOfMeasurementRepositoryImpl) GetQuantityConversion(tx *gorm.DB, payloads masteritempayloads.UomGetQuantityConversion) (masteritempayloads.GetQuantityConversionResponse, *exceptions.BaseErrorResponse) {
	//get uom 1 data base on srouce type and item id
	var quantityResult = payloads.Quantity
	var sellDivided = 1.0

	var result masteritempayloads.GetQuantityConversionResponse
	var itemUomEntities masteritementities.UomItem
	err := tx.Model(&itemUomEntities).Where(masteritementities.UomItem{ItemId: payloads.ItemId, UomSourceTypeCode: payloads.SourceType}).
		First(&itemUomEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("uom data is not found in master uom item"),
			}
		}
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get uom item data",
		}
	}
	if payloads.SourceType == "P" { //purchase
		quantityResult = payloads.Quantity * itemUomEntities.TargetConvertion
		if itemUomEntities.SourceConvertion == 0 {
			itemUomEntities.SourceConvertion = 1
		}
		quantityResult /= itemUomEntities.SourceConvertion
	}
	if payloads.SourceType == "S" {
		if itemUomEntities.SourceUomId == itemUomEntities.TargetUomId {
			sellDivided = 1
		} else {
			sellDivided = itemUomEntities.TargetConvertion
		}
		quantityResult = payloads.Quantity * itemUomEntities.TargetConvertion
		if sellDivided == 0 {
			sellDivided = 1
		}
		quantityResult /= sellDivided
	}
	result.Quantity = payloads.Quantity
	result.QuantityConversion = quantityResult
	result.ItemId = payloads.ItemId
	result.SourceType = payloads.SourceType
	return result, nil

}
