package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type DiscountPercentRepositoryImpl struct {
}

func StartDiscountPercentRepositoryImpl() masteritemrepository.DiscountPercentRepository {
	return &DiscountPercentRepositoryImpl{}
}

// In the GetAllDiscountPercent method
func (r *DiscountPercentRepositoryImpl) GetAllDiscountPercent(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tableStruct := masteritempayloads.DiscountPercentRequest{}
	var orderTypeName string
	newFilterCondition := []utils.FilterCondition{}

	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "order_type_name") {
			orderTypeName = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct).
		Joins("INNER JOIN mtr_discount ON mtr_discount.discount_code_id = mtr_discount_percent.discount_code_id").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_order_type ON mtr_order_type.order_type_id = mtr_discount_percent.order_type_id")

	whereQuery := utils.ApplyFilter(joinTable, newFilterCondition)

	var orderTypeIds []int
	if orderTypeName != "" {
		orderTypeParams := aftersalesserviceapiutils.OrderTypeParams{
			Page: 0, Limit: 100, OrderTypeName: orderTypeName,
		}

		orderTypes, err := aftersalesserviceapiutils.GetAllOrderType(orderTypeParams)
		if err != nil {
			return pages, err
		}

		for _, orderType := range orderTypes {
			orderTypeIds = append(orderTypeIds, orderType.OrderTypeId)
		}

		if len(orderTypeIds) == 0 {
			orderTypeIds = []int{-1}
		}

		whereQuery = whereQuery.Where("mtr_discount_percent.order_type_id IN ?", orderTypeIds)
	}

	var responses []masteritementities.DiscountPercent
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
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

	var mapResponses []masteritempayloads.DiscountPercentListResponse
	for _, response := range responses {
		responseMap := masteritempayloads.DiscountPercentListResponse{
			IsActive:          response.IsActive,
			DiscountPercentId: response.DiscountPercentId,
			DiscountCodeId:    response.DiscountCodeId,
			Discount:          response.Discount,
			OrderTypeId:       response.OrderTypeId,
		}

		var discountDetails masteritempayloads.DiscountDetails
		err := tx.Table("mtr_discount").
			Select("discount_code, discount_description").
			Where("discount_code_id = ?", response.DiscountCodeId).
			Scan(&discountDetails).Error

		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to fetch data from database",
				Err:        err,
			}
		}

		responseMap.DiscountCode = discountDetails.DiscountCode
		responseMap.DiscountDescription = discountDetails.DiscountDescription

		if response.OrderTypeId != 0 {
			orderTypeResponse, err := aftersalesserviceapiutils.GetOrderTypeById(response.OrderTypeId)
			if err != nil {
				responseMap.OrderTypeName = ""
			} else {
				responseMap.OrderTypeName = orderTypeResponse.OrderTypeName
			}
		} else {
			responseMap.OrderTypeName = ""
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
