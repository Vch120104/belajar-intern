package masteritemrepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
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
	var groupId int

	// Extract external filters for groupName and lineTypeCode
	for _, filter := range externalFilter {
		if strings.Contains(filter.ColumnField, "line_type_code") {
			lineTypeCode = filter.ColumnValue
		} else if strings.Contains(filter.ColumnField, "item_group_name") {
			groupName = filter.ColumnValue
		}
	}

	// Filter by item group using GetItemGroupById
	if groupName != "" {
		itemGroupResponse, groupErr := generalserviceapiutils.GetItemGroupById(groupId)
		if groupErr == nil {
			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "item_group_id",
				ColumnValue: strconv.Itoa(itemGroupResponse.ItemGroupId),
			})
		}
	}

	// Filter by line type using GetLineTypeByCode
	if lineTypeCode != "" {
		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeByCode(lineTypeCode)
		if lineErr == nil {
			internalFilter = append(internalFilter, utils.FilterCondition{
				ColumnField: "line_type_id",
				ColumnValue: strconv.Itoa(lineTypeResponse.LineTypeId),
			})
		}
	}

	// Apply internal filters and paginate
	joinTable := utils.CreateJoinSelectStatement(tx, masteritempayloads.ItemClassGetAllResponse{})
	whereQuery := utils.ApplyFilter(joinTable, internalFilter)

	if err := joinTable.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&entities).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch detailed information for item groups and line types
	for i := range entities {
		itemGroupResponse, groupErr := generalserviceapiutils.GetItemGroupById(entities[i].ItemGroupId)
		if groupErr != nil {
			entities[i].ItemGroupName = ""
		} else {
			entities[i].ItemGroupName = itemGroupResponse.ItemGroupName
		}

		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(entities[i].LineTypeId)
		if lineErr != nil {
			entities[i].LineTypeName = ""
		} else {
			entities[i].LineTypeName = lineTypeResponse.LineTypeName
		}
	}

	pages.Rows = entities
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
