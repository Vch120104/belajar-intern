package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type ItemTypeRepositoryImpl struct {
}

func StartItemTypeRepositoryImpl() masteritemrepository.ItemTypeRepository {
	return &ItemTypeRepositoryImpl{}
}

// GetAllItemType fetches the list of item types with pagination and filtering.
func (r *ItemTypeRepositoryImpl) GetAllItemType(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritementities.ItemType

	baseModelQuery := tx.Model(&responses)
	filteredQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	paginatedQuery := filteredQuery.Scopes(pagination.Paginate(&pages, filteredQuery))

	if err := paginatedQuery.Find(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []masteritementities.ItemType{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

// GetItemTypeById fetches a single item type by ID.
func (r *ItemTypeRepositoryImpl) GetItemTypeById(tx *gorm.DB, Id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemType
	var response masteritempayloads.ItemTypeResponse

	if err := tx.Where("item_type_id = ?", Id).First(&entities).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Item type not found",
			Err:        err,
		}
	}

	response = masteritempayloads.ItemTypeResponse{
		ItemTypeId:   entities.ItemTypeId,
		ItemTypeCode: entities.ItemTypeCode,
		ItemTypeName: entities.ItemTypeName,
		IsActive:     entities.IsActive,
	}

	return response, nil
}

// GetItemTypeByCode fetches a single item type by its code.
func (r *ItemTypeRepositoryImpl) GetItemTypeByCode(tx *gorm.DB, code string) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemType
	var response masteritempayloads.ItemTypeResponse

	if err := tx.Where("item_type_code = ?", code).First(&entities).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response = masteritempayloads.ItemTypeResponse{
		ItemTypeId:   entities.ItemTypeId,
		ItemTypeCode: entities.ItemTypeCode,
		ItemTypeName: entities.ItemTypeName,
		IsActive:     entities.IsActive,
	}

	return response, nil
}

// CreateItemType creates a new item type.
func (r *ItemTypeRepositoryImpl) CreateItemType(tx *gorm.DB, request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse) {

	if request.ItemTypeCode == "" || request.ItemTypeName == "" {
		return masteritementities.ItemType{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "ItemTypeCode and ItemTypeName are required",
		}
	}

	entities := masteritementities.ItemType{
		ItemTypeCode: request.ItemTypeCode,
		ItemTypeName: request.ItemTypeName,
	}

	if err := tx.Create(&entities).Error; err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create item type",
			Err:        err,
		}
	}

	return entities, nil
}

// SaveItemType saves a new or updates an existing item type.
func (r *ItemTypeRepositoryImpl) SaveItemType(tx *gorm.DB, id int, request masteritempayloads.ItemTypeRequest) (masteritementities.ItemType, *exceptions.BaseErrorResponse) {

	if request.ItemTypeCode == "" || request.ItemTypeName == "" {
		return masteritementities.ItemType{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "ItemTypeCode and ItemTypeName are required",
		}
	}

	var entities masteritementities.ItemType
	if id != 0 {
		if err := tx.Where("item_type_id = ?", id).First(&entities).Error; err != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "ItemType not found",
				Err:        err,
			}
		}
	}

	entities.ItemTypeCode = request.ItemTypeCode
	entities.ItemTypeName = request.ItemTypeName

	if err := tx.Save(&entities).Error; err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save item type",
			Err:        err,
		}
	}

	return entities, nil
}

// ChangeStatusItemType toggles the active status of an item type.
func (r *ItemTypeRepositoryImpl) ChangeStatusItemType(tx *gorm.DB, Id int) (masteritempayloads.ItemTypeResponse, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemType
	var response masteritempayloads.ItemTypeResponse

	if err := tx.Where("item_type_id = ?", Id).First(&entities).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entities.IsActive = !entities.IsActive

	if err := tx.Save(&entities).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to change item type status",
			Err:        err,
		}
	}

	response.ItemTypeId = entities.ItemTypeId
	response.ItemTypeCode = entities.ItemTypeCode
	response.ItemTypeName = entities.ItemTypeName
	response.IsActive = entities.IsActive

	return response, nil
}

// GetItemTypeDropDown fetches a list of item types for a dropdown.
func (r *ItemTypeRepositoryImpl) GetItemTypeDropDown(tx *gorm.DB) ([]masteritempayloads.ItemTypeDropDownResponse, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemTypeDropDownResponse

	if err := tx.Model(&masteritementities.ItemType{}).
		Select("item_type_id, item_type_code, item_type_name, is_active").
		Find(&responses).Error; err != nil {

		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve item type dropdown data",
			Err:        err,
		}
	}

	return responses, nil
}
