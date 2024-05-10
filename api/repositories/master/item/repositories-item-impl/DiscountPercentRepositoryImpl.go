package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type DiscountPercentRepositoryImpl struct {
}

func StartDiscountPercentRepositoryImpl() masteritemrepository.DiscountPercentRepository {
	return &DiscountPercentRepositoryImpl{}
}

func (r *DiscountPercentRepositoryImpl) GetAllDiscountPercent(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.DiscountPercentListResponse
	var getOrderTypeResponse []masteritempayloads.OrderTypeResponse

	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var orderTypeName string
	responseStruct := reflect.TypeOf(masteritempayloads.DiscountPercentListResponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
		}
	}

	//apply external services filter
	for i := 0; i < len(externalServiceFilter); i++ {
		orderTypeName = externalServiceFilter[i].ColumnValue
	}

	// define table struct
	tableStruct := masteritempayloads.DiscountPercentListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Execute the query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	// Scan the results into the responses slice
	for rows.Next() {
		var response masteritempayloads.DiscountPercentListResponse
		if err := rows.Scan(&response.IsActive, &response.DiscountPercentId, &response.DiscountCodeId, &response.DiscountCodeValue, &response.DiscountCodeDescription, &response.OrderTypeId, &response.Discount); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	// Fetch order type data
	orderTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/order-type-filter?order_type_name=" + orderTypeName
	errUrlDiscountPercent := utils.Get(orderTypeUrl, &getOrderTypeResponse, nil)
	if errUrlDiscountPercent != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlDiscountPercent,
		}
	}

	// Perform inner join with order type data
	joinedData := utils.DataFrameInnerJoin(responses, getOrderTypeResponse, "OrderTypeId")

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *DiscountPercentRepositoryImpl) GetDiscountPercentById(tx *gorm.DB, Id int) (masteritempayloads.DiscountPercentResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.DiscountPercent{}
	response := masteritempayloads.DiscountPercentResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.DiscountPercent{
			DiscountPercentId: Id,
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

func (r *DiscountPercentRepositoryImpl) SaveDiscountPercent(tx *gorm.DB, request masteritempayloads.DiscountPercentResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
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

func (r *DiscountPercentRepositoryImpl) ChangeStatusDiscountPercent(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.DiscountPercent

	result := tx.Model(&entities).
		Where("discount_percent_id = ?", Id).
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
