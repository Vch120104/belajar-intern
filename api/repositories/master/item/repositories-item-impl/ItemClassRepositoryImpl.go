package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
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
	if err := tx.Model(&entities).Where(masteritementities.ItemClass{ItemGroupID: groupId}).Scan(&response).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

// GetItemClassByCode implements masteritemrepository.ItemClassRepository.
func (r *ItemClassRepositoryImpl) GetItemClassByCode(tx *gorm.DB, itemClassCode string) (masteritempayloads.ItemClassResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemClass{}
	response := masteritempayloads.ItemClassResponse{}

	err := tx.Model(&entities).
		Select("is_active", "item_class_id", "item_class_code", "item_class_name", "item_group_id", "line_type_id").
		Where(masteritementities.ItemClass{ItemClassCode: itemClassCode}).First(&response).Error
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

	if response.LineTypeId != 0 {
		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(response.LineTypeId)
		if lineErr != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: lineErr.StatusCode,
				Message:    "Error fetching line type data",
				Err:        lineErr.Err,
			}
		}

		response.LineTypeCode = lineTypeResponse.LineTypeCode
		response.LineTypeName = lineTypeResponse.LineTypeName
	} else {
		response.LineTypeName = ""
		response.LineTypeCode = ""
	}

	if response.ItemGroupId != 0 {
		var itemGroupResponse masteritementities.ItemGroup
		if err := tx.Where("item_group_id = ?", response.ItemGroupId).First(&itemGroupResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with id %d not found", response.ItemGroupId),
				}
			}
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group ",
				Err:        err,
			}
		}

		response.ItemGroupName = itemGroupResponse.ItemGroupName
	} else {
		response.ItemGroupName = ""
	}

	return response, nil
}

// GetItemClassDropDown implements masteritemrepository.ItemClassRepository.
func (r *ItemClassRepositoryImpl) GetItemClassDropDown(tx *gorm.DB) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	entities := []masteritementities.ItemClass{}
	response := []masteritempayloads.ItemClassDropdownResponse{}
	if err := tx.Model(&entities).Scan(&response).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *ItemClassRepositoryImpl) GetItemClassMfgDropdown(tx *gorm.DB) ([]masteritempayloads.ItemClassDropdownResponse, *exceptions.BaseErrorResponse) {
	response := []masteritempayloads.ItemClassDropdownResponse{}

	err := tx.Table("mtr_item_class").
		Select(`is_active, item_class_id, item_class_name`).
		Where("is_manufacturing_item_type = 1 AND is_active = 1").
		Order("item_class_name").
		Scan(&response).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching item class record",
			Err:        err,
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
		var itemGroupResponse masteritementities.ItemGroup
		err := tx.Where("item_group_id = ?", groupId).First(&itemGroupResponse).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item group not found",
					Err:        fmt.Errorf("item group with id %d not found", groupId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group",
				Err:        err,
			}
		}

		internalFilter = append(internalFilter, utils.FilterCondition{
			ColumnField: "item_group_id",
			ColumnValue: strconv.Itoa(itemGroupResponse.ItemGroupId),
		})
	}

	// Apply internal filters and paginate
	joinTable := utils.CreateJoinSelectStatement(tx, masteritempayloads.ItemClassGetAllResponse{})
	whereQuery := utils.ApplyFilter(joinTable, internalFilter)

	// Filter by line type using GetLineTypeByCode
	if lineTypeCode != "" {
		lineTypeParam := generalserviceapiutils.LineTypeListParams{
			Page:         0,
			Limit:        1000,
			LineTypeCode: lineTypeCode,
		}

		lineTypeResponse, errLine := generalserviceapiutils.GetLineTypeListByCode(lineTypeParam)
		if errLine != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errLine.StatusCode,
				Message:    "Error fetching supplier data",
				Err:        errLine.Err,
			}
		}

		var lineTypeIds []int
		for _, lineType := range lineTypeResponse {
			lineTypeIds = append(lineTypeIds, lineType.LineTypeId)
		}
		if len(lineTypeIds) != 0 {
			whereQuery = whereQuery.Where("line_type_id IN ?", lineTypeIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	if err := joinTable.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&entities).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch detailed information for item groups and line types
	for i := range entities {
		var itemGroupResponse masteritementities.ItemGroup
		err := tx.Where("item_group_id = ?", entities[i].ItemGroupId).First(&itemGroupResponse).Error
		if err != nil {

			entities[i].ItemGroupName = ""

			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item group",
				Err:        err,
			}
		}

		entities[i].ItemGroupName = itemGroupResponse.ItemGroupName

		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(entities[i].LineTypeId)
		if lineErr != nil {
			entities[i].LineTypeName = ""
			entities[i].LineTypeCode = ""
		} else {
			entities[i].LineTypeName = lineTypeResponse.LineTypeName
			entities[i].LineTypeCode = lineTypeResponse.LineTypeCode
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

	if response.LineTypeId != 0 {

		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(response.LineTypeId)
		if lineErr != nil {

			return response, &exceptions.BaseErrorResponse{
				StatusCode: lineErr.StatusCode,
				Err:        lineErr.Err,
			}
		}

		response.LineTypeName = lineTypeResponse.LineTypeName
	}

	return response, nil
}

func (r *ItemClassRepositoryImpl) SaveItemClass(tx *gorm.DB, request masteritempayloads.ItemClassResponse) (bool, *exceptions.BaseErrorResponse) {

	var itemGroup masteritementities.ItemGroup
	if err := tx.Where("item_group_id = ?", request.ItemGroupId).First(&itemGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item group not found",
				Err:        fmt.Errorf("item group with id %d not found", request.ItemGroupId),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item group code",
			Err:        err,
		}
	}

	// CHECK LINE TYPE ID IF ITEM GROUP IS 'INVENTORY'
	if itemGroup.ItemGroupName == "Inventory" || itemGroup.ItemGroupCode == "IN" {

		lineType, lineTypeErr := generalserviceapiutils.GetLineTypeById(request.LineTypeId)
		if lineTypeErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: lineTypeErr.StatusCode,
				Err:        lineTypeErr.Err,
			}
		}

		if lineType == (generalserviceapiutils.LineTypeResponse{}) {
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
