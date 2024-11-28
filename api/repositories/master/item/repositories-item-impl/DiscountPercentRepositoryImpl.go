package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type DiscountPercentRepositoryImpl struct {
}

func StartDiscountPercentRepositoryImpl() masteritemrepository.DiscountPercentRepository {
	return &DiscountPercentRepositoryImpl{}
}

func (r *DiscountPercentRepositoryImpl) GetAllDiscountPercent(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tableStruct := masteritempayloads.DiscountPercentRequest{}
	var orderTypeName string
	newFilterCondition := []utils.FilterCondition{}

	// Separate order type filters and other conditions
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "order_type_name") {
			orderTypeName = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	// Build the join table query
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct).
		Joins("LEFT JOIN mtr_discount ON mtr_discount.discount_code_id = mtr_discount_percent.discount_code_id").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_order_type ON mtr_order_type.order_type_id = mtr_discount_percent.order_type_id")

	// Apply filter conditions to the query
	whereQuery := utils.ApplyFilter(joinTable, newFilterCondition)

	// Handle order_type_name filter via external service
	var orderTypeIds []int
	if orderTypeName != "" {
		orderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type?page=0&limit=100&order_type_name=" + orderTypeName
		var getOrderTypeResponse []masteritempayloads.OrderTypeResponse

		if err := utils.Get(orderTypeURL, &getOrderTypeResponse, nil); err == nil {
			for _, orderType := range getOrderTypeResponse {
				orderTypeIds = append(orderTypeIds, orderType.OrderTypeId)
			}
		}

		if len(orderTypeIds) == 0 {
			orderTypeIds = []int{-1} // Ensure no matches
		}

		whereQuery = whereQuery.Where("mtr_discount_percent.order_type_id IN ?", orderTypeIds)
	}

	// Execute the query with pagination
	var responses []masteritempayloads.DiscountPercentListResponse
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch data from database",
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	// Prepare the response map
	var mapResponses []map[string]interface{}
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"is_active":            response.IsActive,
			"discount_percent_id":  response.DiscountPercentId,
			"discount_code_id":     response.DiscountCodeId,
			"discount_code":        response.DiscountCode,
			"discount_description": response.DiscountDescription,
			"order_type_id":        response.OrderTypeId,
			"discount":             response.Discount,
		}

		// Fetch additional order type details if applicable
		if response.OrderTypeId != 0 {
			orderTypeURL := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(response.OrderTypeId)
			var getOrderTypeResponse masteritempayloads.OrderTypeResponse

			if err := utils.Get(orderTypeURL, &getOrderTypeResponse, nil); err == nil {
				responseMap["order_type_name"] = getOrderTypeResponse.OrderTypeName
			} else {
				responseMap["order_type_name"] = ""
			}
		} else {
			responseMap["order_type_name"] = ""
		}

		mapResponses = append(mapResponses, responseMap)
	}

	pages.Rows = mapResponses
	return pages, nil
}

func (r *DiscountPercentRepositoryImpl) GetDiscountPercentById(tx *gorm.DB, Id int) (masteritempayloads.DiscountPercentResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.DiscountPercent{}
	response := masteritempayloads.DiscountPercentResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.DiscountPercent{
			DiscountPercentId: Id,
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

func (r *DiscountPercentRepositoryImpl) SaveDiscountPercent(tx *gorm.DB, request masteritempayloads.DiscountPercentResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.DiscountPercent{
		IsActive:          request.IsActive,
		DiscountPercentId: request.DiscountPercentId,
		DiscountCodeId:    request.DiscountCodeId,
		OrderTypeId:       request.OrderTypeId,
		Discount:          request.Discount,
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

func (r *DiscountPercentRepositoryImpl) ChangeStatusDiscountPercent(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.DiscountPercent

	result := tx.Model(&entities).
		Where("discount_percent_id = ?", Id).
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
