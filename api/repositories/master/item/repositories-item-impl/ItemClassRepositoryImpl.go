package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type ItemClassRepositoryImpl struct {
}

func StartItemClassRepositoryImpl() masteritemrepository.ItemClassRepository {
	return &ItemClassRepositoryImpl{}
}

// GetItemClassDropDownbyGroupId implements masteritemrepository.ItemClassRepository.
func (r *ItemClassRepositoryImpl) GetItemClassDropDownbyGroupId(tx *gorm.DB, groupId int) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemClass{}
	response := []masteritempayloads.ItemClassDropdownResponse{}
	if err := tx.Model(entities).Where(masteritementities.ItemClass{ItemGroupID: groupId}).Scan(&response).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(response) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(""),
		}
	}
	return response, nil
}

// GetItemClassByCode implements masteritemrepository.ItemClassRepository.
func (r *ItemClassRepositoryImpl) GetItemClassByCode(tx *gorm.DB, itemClassCode string) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemClass{}
	response := masteritempayloads.ItemClassResponse{}

	err := tx.Model(&entities).Select("mtr_item_class.*").
		Where(masteritementities.ItemClass{
			ItemClassCode: itemClassCode,
		}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	lineTypeResponse := masteritempayloads.LineTypeResponse{}

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(response.LineTypeId)

	if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData, errdf := utils.DataFrameInnerJoin([]masteritempayloads.ItemClassResponse{response}, []masteritempayloads.LineTypeResponse{lineTypeResponse}, "LineTypeId")

	if errdf != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	value, ok := joinedData[0]["LineTypeName_1"]

	if ok {
		switch v := value.(type) {
		case string:
			response.LineTypeName = v
		}
	}

	return response, nil
}

// GetItemClassDropDown implements masteritemrepository.ItemClassRepository.
func (r *ItemClassRepositoryImpl) GetItemClassDropDown(tx *gorm.DB) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemClass{}
	response := []masteritempayloads.ItemClassDropdownResponse{}
	if err := tx.Model(entities).Scan(&response).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(response) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(""),
		}
	}
	return response, nil
}

func (r *ItemClassRepositoryImpl) GetAllItemClass(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemClass{}
	var responses []masteritempayloads.ItemClassGetAllResponse
	var getLineTypeResponse []masteritempayloads.LineTypeResponse
	var getItemGroupResponse []masteritempayloads.ItemGroupResponse
	var groupName, lineTypeCode string

	//apply external services filter
	for _, value := range externalFilter {
		if strings.Contains(value.ColumnField, "line_type_code") {
			lineTypeCode = value.ColumnValue
		} else {
			groupName = value.ColumnValue
		}
	}

	//IF GROUP NAME ON FILTER (EXTERNAL)
	if groupName != "" {
		groupServiceUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_name=" + groupName

		errUrlItemGroup := utils.Get(groupServiceUrl, &getItemGroupResponse, nil)
		if errUrlItemGroup != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlItemGroup,
			}
		}

		if len(getItemGroupResponse) == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNoContent,
				Err:        errors.New("item group not found"),
			}
		}

		var ids string
		for _, value := range getItemGroupResponse {
			ids += fmt.Sprintf("%d,", value.ItemGroupId)
		}

		condition := utils.FilterCondition{
			//add tag -> #multiple, to initiate filter for multiple ID
			ColumnField: "item_group_id #multiple",
			ColumnValue: strings.TrimSuffix(ids, ","),
		}

		internalFilter = append(internalFilter, condition)
	}

	//IF LINE TYPE ON FILTER
	if lineTypeCode != "" {
		lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type?line_type_code=" + lineTypeCode

		errUrlLineType := utils.Get(lineTypeUrl, &getLineTypeResponse, nil)

		if errUrlLineType != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlLineType,
			}
		}

		if len(getLineTypeResponse) == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNoContent,
				Err:        errors.New("line type not found"),
			}
		}

		var ids string
		for _, value := range getLineTypeResponse {
			ids += fmt.Sprintf("%d,", value.LineTypeId)
		}

		condition := utils.FilterCondition{
			//add tag -> #multiple, to initiate filter for multiple ID
			ColumnField: "line_type_id #multiple",
			ColumnValue: strings.TrimSuffix(ids, ","),
		}

		internalFilter = append(internalFilter, condition)
	}

	// fmt.Println("internal filter", internalServiceFilter)
	// fmt.Println("external filter", externalServiceFilter)

	//define base model
	baseModelQuery := tx.Model(&entities)
	//apply where query
	whereQuery := utils.ApplyFilter(baseModelQuery, internalFilter)
	//apply pagination and execute
	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New("no contents"),
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() *exceptions.BaseErrorResponse {
		defer wg.Done()
		groupServiceUrl := config.EnvConfigs.GeneralServiceUrl + "filter-item-group?item_group_name=" + groupName

		errUrlItemGroup := utils.Get(groupServiceUrl, &getItemGroupResponse, nil)

		if errUrlItemGroup != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errUrlItemGroup,
			}
		}

		return nil
	}()

	wg.Add(1)
	go func() *exceptions.BaseErrorResponse {
		defer wg.Done()

		lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type?line_type_code=" + lineTypeCode

		errUrlLineType := utils.Get(lineTypeUrl, &getLineTypeResponse, nil)

		if errUrlLineType != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errUrlLineType,
			}
		}
		return nil
	}()

	wg.Wait()

	joinedData, errdf := utils.DataFrameInnerJoin(responses, getItemGroupResponse, "ItemGroupId")

	if errdf != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	joinedDataSecond, errdf := utils.DataFrameInnerJoin(joinedData, getLineTypeResponse, "LineTypeId")

	if errdf != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	pages.Rows = joinedDataSecond

	return pages, nil

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

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(response.LineTypeId)
	if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData, errdf := utils.DataFrameInnerJoin([]masteritempayloads.ItemClassResponse{response}, []masteritempayloads.LineTypeResponse{lineTypeResponse}, "LineTypeId")

	if errdf != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

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
	var getLineTypeResponse masteritempayloads.LineTypeResponse
	var getItemGroupResponse masteritempayloads.ItemGroupResponse

	//CHECK ITEM GROUP ID
	groupUrl := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(request.ItemGroupId)

	errUrlItemGroup := utils.Get(groupUrl, &getItemGroupResponse, nil)

	if errUrlItemGroup != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlItemGroup,
		}
	}

	if getItemGroupResponse == (masteritempayloads.ItemGroupResponse{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("item group not found"),
		}
	}

	//CHECK LINE TYPE ID IF ITEM GROUP IS 'INVENTORY'
	if getItemGroupResponse.ItemGroupName == "Inventory" || getItemGroupResponse.ItemGroupCode == "IN" {
		lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(request.LineTypeId)
		errUrlLineType := utils.Get(lineTypeUrl, &getLineTypeResponse, nil)

		if errUrlLineType != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errUrlLineType,
			}
		}

		if getLineTypeResponse == (masteritempayloads.LineTypeResponse{}) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("line type not found"),
			}
		}
	} else {
		request.LineTypeId = 0
	}

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
