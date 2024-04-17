package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type ItemLocationRepositoryImpl struct {
}

func StartItemLocationRepositoryImpl() masteritemrepository.ItemLocationRepository {
	return &ItemLocationRepositoryImpl{}
}

func (r *ItemLocationRepositoryImpl) GetAllItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	entities := []masteritementities.ItemLocation{}
	var responses []masteritempayloads.ItemLocationResponse
	var getWarehouseGroupResponse []masteritempayloads.ItemLocWarehouseGroupResponse
	var getItemResponse []masteritempayloads.ItemLocResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var groupName, lineTypeCode string
	responseStruct := reflect.TypeOf(masteritempayloads.ItemLocationResponse{})

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
		if strings.Contains(externalServiceFilter[i].ColumnField, "warehouse_group_id") {
			lineTypeCode = externalServiceFilter[i].ColumnValue
		} else {
			groupName = externalServiceFilter[i].ColumnValue
		}
	}

	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	groupServiceUrl := "http://10.1.32.26:8000/general-service/api/general/filter-item-group?item_group_name=" + groupName

	errUrlItemGroup := utils.Get(groupServiceUrl, &getWarehouseGroupResponse, nil)

	if errUrlItemGroup != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlItemGroup,
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemResponse, "ItemGroupId")

	lineTypeUrl := "http://10.1.32.26:8000/general-service/api/general/line-type?line_type_code=" + lineTypeCode

	errUrlLineType := utils.Get(lineTypeUrl, &getItemResponse, nil)

	if errUrlLineType != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlLineType,
		}
	}

	joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getWarehouseGroupResponse, "LineTypeId")

	return joinedDataSecond, len(joinedDataSecond), len(joinedDataSecond), nil
}

func (r *ItemLocationRepositoryImpl) SaveItemLocation(tx *gorm.DB, request masteritempayloads.ItemLocationRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{
		WarehouseGroupId:   request.WarehouseGroupId,
		WarehouseGroupCode: request.WarehouseGroupCode,
		ItemId:             request.ItemId,
		ItemCode:           request.ItemCode,
		ItemName:           request.ItemName,
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
