package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
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

func (r *MarkupRateRepositoryImpl) GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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

	// Execute the query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	// Initialize responses slice with 0 length
	responses = make([]masteritempayloads.MarkupRateListResponse, 0)

	// Scan the results into the responses slice
	for rows.Next() {
		var response masteritempayloads.MarkupRateListResponse
		if err := rows.Scan(&response.IsActive, &response.MarkupRateId, &response.MarkupMasterId, &response.MarkupMasterCode, &response.MarkupMasterDescription, &response.OrderTypeId, &response.MarkupRate); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	// Fetch order type data
	orderTypeUrl := config.EnvConfigs.AfterSalesServiceUrl + "order-type/" + orderTypeName
	errUrlMarkupRate := utils.Get(orderTypeUrl, &getOrderTypeResponse, nil)
	if errUrlMarkupRate != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlMarkupRate,
		}
	}

	// Perform inner join with order type data
	joinedData, errdf := utils.DataFrameInnerJoin(responses, getOrderTypeResponse, "OrderTypeId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
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
