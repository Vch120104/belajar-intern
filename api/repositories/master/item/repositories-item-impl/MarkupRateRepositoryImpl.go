package masteritemrepositoryimpl

import (
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

type MarkupRateRepositoryImpl struct {
}

func StartMarkupRateRepositoryImpl() masteritemrepository.MarkupRateRepository {
	return &MarkupRateRepositoryImpl{}
}

func (r *MarkupRateRepositoryImpl) GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.MarkupRateListResponse
	var getOrderTypeResponse []masteritempayloads.OrderTypeResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var orderTypeName string
	responseStruct := reflect.TypeOf(masteritempayloads.MarkupRateListResponse{})

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
	tableStruct := masteritempayloads.MarkupRateListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)
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
			Err:        errors.New(""),
		}
	}

	defer rows.Close()

	orderTypeUrl := "http://10.1.32.26:8000/general-service/api/general/order-type-filter?order_type_name=" + orderTypeName

	errUrlMarkupRate := utils.Get(orderTypeUrl, &getOrderTypeResponse, nil)

	if errUrlMarkupRate != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlMarkupRate,
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getOrderTypeResponse, "OrderTypeId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *MarkupRateRepositoryImpl) GetMarkupRateById(tx *gorm.DB, Id int) (masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.MarkupRate{}
	response := masteritempayloads.MarkupRateResponse{}

	rows, err := tx.Model(&entities).
		Where(masteritementities.MarkupRate{
			MarkupRateId: Id,
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

func (r *MarkupRateRepositoryImpl) SaveMarkupRate(tx *gorm.DB, request masteritempayloads.MarkupRateRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
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

func (r *MarkupRateRepositoryImpl) ChangeStatusMarkupRate(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteritementities.MarkupRate

	result := tx.Model(&entities).
		Where("markup_rate_id = ?", Id).
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
