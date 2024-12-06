package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type MarkupRateRepositoryImpl struct {
}

func StartMarkupRateRepositoryImpl() masteritemrepository.MarkupRateRepository {
	return &MarkupRateRepositoryImpl{}
}

func (r *MarkupRateRepositoryImpl) GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []masteritempayloads.MarkupRateListResponse
	var orderTypeName string
	newFilterCondition := []utils.FilterCondition{}

	// Processing filter conditions
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "order_type_name") {
			orderTypeName = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	// Base query
	baseQuery := tx.Model(&masteritementities.MarkupRate{}).
		Joins("INNER JOIN mtr_markup_master ON mtr_markup_master.markup_master_id = mtr_markup_rate.markup_master_id").
		Joins("LEFT JOIN mtr_order_type ON mtr_order_type.order_type_id = mtr_markup_rate.order_type_id").
		Select("mtr_markup_rate.*, mtr_markup_master.markup_code, mtr_markup_master.markup_description")

	// Apply filter conditions
	whereQuery := utils.ApplyFilter(baseQuery, newFilterCondition)

	// Handle order type filters
	var orderTypeIds []int
	if orderTypeName != "" {
		orderTypeParams := aftersalesserviceapiutils.OrderTypeParams{
			Page:          0,
			Limit:         100,
			OrderTypeName: orderTypeName,
		}
		orderTypes, err := aftersalesserviceapiutils.GetAllOrderType(orderTypeParams)
		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: err.StatusCode,
				Message:    "Failed to fetch order types",
				Err:        err.Err,
			}
		}

		for _, orderType := range orderTypes {
			orderTypeIds = append(orderTypeIds, orderType.OrderTypeId)
		}
		if len(orderTypeIds) == 0 {
			orderTypeIds = []int{-1}
		}

		whereQuery = whereQuery.Where("mtr_markup_rate.order_type_id IN ?", orderTypeIds)
	}

	// Fetching data from the database
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data from the database",
			Err:        err,
		}
	}

	// If no results
	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	// Prepare the response
	var responses []map[string]interface{}
	for _, entity := range entities {
		response := map[string]interface{}{
			"is_active":          entity.IsActive,
			"markup_rate_id":     entity.MarkupRateId,
			"markup_master_id":   entity.MarkupMasterId,
			"order_type_id":      entity.OrderTypeId,
			"markup_rate":        entity.MarkupRate,
			"markup_code":        entity.MarkupCode,
			"markup_description": entity.MarkupDescription,
		}

		// If order type exists, fetch the order type name
		if entity.OrderTypeId != 0 {
			orderType, err := aftersalesserviceapiutils.GetOrderTypeById(entity.OrderTypeId)
			if err == nil {
				response["order_type_name"] = orderType.OrderTypeName
			} else {
				response["order_type_name"] = ""
			}
		}
		responses = append(responses, response)
	}
	fmt.Println(entities)
	pages.Rows = responses
	return pages, nil
}

func (r *MarkupRateRepositoryImpl) GetMarkupRateById(tx *gorm.DB, Id int) (masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.MarkupRate{}
	response := masteritempayloads.MarkupRateResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.MarkupRate{
			MarkupRateId: Id,
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

func (r *MarkupRateRepositoryImpl) SaveMarkupRate(tx *gorm.DB, request masteritempayloads.MarkupRateRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.MarkupRate{
		IsActive:       true,
		MarkupRateId:   request.MarkupRateId,
		MarkupMasterId: request.MarkupMasterId,
		OrderTypeId:    request.OrderTypeId,
		MarkupRate:     request.MarkupRate,
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

func (r *MarkupRateRepositoryImpl) ChangeStatusMarkupRate(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.MarkupRate

	result := tx.Model(&entities).
		Where("markup_rate_id = ?", Id).
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

func (r *MarkupRateRepositoryImpl) GetMarkupRateByMarkupMasterAndOrderType(tx *gorm.DB, MarkupMasterId int, OrderTypeId int) ([]masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.MarkupRate{}
	response := []masteritempayloads.MarkupRateResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.MarkupRate{
			MarkupMasterId: MarkupMasterId,
			OrderTypeId:    OrderTypeId,
		}).
		Find(&response).
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
