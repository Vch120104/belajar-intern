package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemClassRepositoryImpl struct {
}

func StartItemClassRepositoryImpl() masteritemrepository.ItemClassRepository {
	return &ItemClassRepositoryImpl{}
}

func (r *ItemClassRepositoryImpl) GetAllItemClass(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemClass{}
	var responses []masteritempayloads.ItemClassResponse
	var getLineTypeResponse []masteritempayloads.LineTypeResponse
	var getItemGroupResponse []masteritempayloads.ItemGroupResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var groupName, lineTypeCode string
	responseStruct := reflect.TypeOf(masteritempayloads.ItemClassResponse{})

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
		if strings.Contains(externalServiceFilter[i].ColumnField, "line_type_code") {
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
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	groupServiceUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_name=" + groupName

	errUrlItemGroup := utils.Get(groupServiceUrl, &getItemGroupResponse, nil)

	if errUrlItemGroup != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlItemGroup,
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type?line_type_code=" + lineTypeCode

	errUrlLineType := utils.Get(lineTypeUrl, &getLineTypeResponse, nil)

	if errUrlLineType != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlLineType,
		}
	}

	joinedDataSecond := utils.DataFrameInnerJoin(joinedData, getLineTypeResponse, "LineTypeId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedDataSecond, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemClassRepositoryImpl) GetItemClassById(tx *gorm.DB, Id int) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemClass{}
	response := masteritempayloads.ItemClassResponse{}

	err := tx.Model(&entities).Select("mtr_item_class.*").
		Where(masteritementities.ItemClass{
			ItemClassId: Id,
		}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	lineTypeResponse := masteritempayloads.LineTypeResponse{}

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/line-type/" + strconv.Itoa(response.LineTypeId)

	if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData := utils.DataFrameInnerJoin([]masteritempayloads.ItemClassResponse{response}, []masteritempayloads.LineTypeResponse{lineTypeResponse}, "LineTypeId")

	value, ok := joinedData[0]["LineTypeName_1"]

	if ok {
		switch v := value.(type) {
		case string:
			response.LineTypeName = v
		}
	}

	return response, nil
}

func (r *ItemClassRepositoryImpl) SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemClass{
		IsActive:      request.IsActive,
		ItemClassId:   request.ItemClassId,
		ItemClassCode: request.ItemClassCode,
		ItemGroupID:   request.ItemGroupId,
		LineTypeID:    request.LineTypeId,
		ItemClassName: request.ItemClassName,
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

func (r *ItemClassRepositoryImpl) ChangeStatusItemClass(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemClass

	result := tx.Model(&entities).
		Where("item_class_id = ?", Id).
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
