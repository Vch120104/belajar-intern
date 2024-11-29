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
	entities := []masteritempayloads.ItemClassResponse{}
	var groupName, lineTypeCode string

	for _, filter := range externalFilter {
		if strings.Contains(filter.ColumnField, "line_type_code") {
			lineTypeCode = filter.ColumnValue
		} else if strings.Contains(filter.ColumnField, "item_group_name") {
			groupName = filter.ColumnValue
		}
	}

	// Filter berdasarkan group name
	if groupName != "" {
		groupServiceURL := config.EnvConfigs.GeneralServiceUrl + "item-group?page=0&limit=100&item_group_name=" + groupName
		var itemGroups []masteritempayloads.ItemGroupResponse

		if err := utils.Get(groupServiceURL, &itemGroups, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if len(itemGroups) == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNoContent,
				Err:        errors.New("item group not found"),
			}
		}

		var ids string
		for _, group := range itemGroups {
			ids += fmt.Sprintf("%d,", group.ItemGroupId)
		}

		internalFilter = append(internalFilter, utils.FilterCondition{
			ColumnField: "item_group_id #multiple",
			ColumnValue: strings.TrimSuffix(ids, ","),
		})
	}

	// Filter berdasarkan line type
	if lineTypeCode != "" {
		lineTypeURL := config.EnvConfigs.GeneralServiceUrl + "line-types?page=0&limit=100&line_type_code=" + lineTypeCode
		var lineTypes []masteritempayloads.LineTypeResponse

		if err := utils.Get(lineTypeURL, &lineTypes, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if len(lineTypes) == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNoContent,
				Err:        errors.New("line type not found"),
			}
		}

		var ids string
		for _, lineType := range lineTypes {
			ids += fmt.Sprintf("%d,", lineType.LineTypeId)
		}

		internalFilter = append(internalFilter, utils.FilterCondition{
			ColumnField: "line_type_id #multiple",
			ColumnValue: strings.TrimSuffix(ids, ","),
		})
	}

	joinTable := utils.CreateJoinSelectStatement(tx, masteritempayloads.ItemClassGetAllResponse{})
	whereQuery := utils.ApplyFilter(joinTable, internalFilter)

	if err := joinTable.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&entities).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var itemGroupIds []int
	var lineTypeIds []int

	for _, entity := range entities {
		if entity.ItemGroupId != 0 {
			itemGroupIds = append(itemGroupIds, entity.ItemGroupId)
		}
		if entity.LineTypeId != 0 {
			lineTypeIds = append(lineTypeIds, entity.LineTypeId)
		}
	}

	// Get item group names based on itemGroupIds
	if len(itemGroupIds) > 0 {
		groupServiceURL := fmt.Sprintf("%sitem-group-multi-id/%s", config.EnvConfigs.GeneralServiceUrl, strings.Join(toStringList(itemGroupIds), ","))
		var itemGroups []masteritempayloads.ItemGroupResponse

		if err := utils.Get(groupServiceURL, &itemGroups, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		itemGroupNames := map[int]string{}
		for _, group := range itemGroups {
			itemGroupNames[group.ItemGroupId] = group.ItemGroupName
		}

		for i := range entities {
			if groupName, found := itemGroupNames[entities[i].ItemGroupId]; found {
				entities[i].ItemGroupName = groupName
			}
		}
	}

	// Get line type names based on lineTypeIds
	if len(lineTypeIds) > 0 {
		lineTypeURL := fmt.Sprintf("%sline-type-list?page=0&limit=100&line_type_ids=%s", config.EnvConfigs.GeneralServiceUrl, strings.Join(toStringList(lineTypeIds), ","))
		var lineTypes []masteritempayloads.LineTypeResponse

		if err := utils.Get(lineTypeURL, &lineTypes, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		lineTypeNames := map[int]string{}
		for _, lineType := range lineTypes {
			lineTypeNames[lineType.LineTypeId] = lineType.LineTypeName
		}

		for i := range entities {
			if lineTypeName, found := lineTypeNames[entities[i].LineTypeId]; found {
				entities[i].LineTypeName = lineTypeName
			}
		}
	}

	pages.Rows = entities

	return pages, nil
}

// Helper function to convert a list of integers to a list of strings
func toStringList(ids []int) []string {
	var result []string
	for _, id := range ids {
		result = append(result, fmt.Sprintf("%d", id))
	}
	return result
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item class not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	lineTypeResponse := masteritempayloads.LineTypeResponse{}
	if response.LineTypeId != 0 {

		lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(response.LineTypeId)
		if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
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
