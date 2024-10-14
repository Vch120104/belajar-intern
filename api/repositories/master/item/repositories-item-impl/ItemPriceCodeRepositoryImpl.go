package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemPriceCodeRepositoryImpl struct {
}

func StartItemPriceCodeImpl() masteritemlevelrepo.ItemPriceCodeRepository {
	return &ItemPriceCodeRepositoryImpl{}
}

// GetAllItemPriceCode implements masteritemrepository.ItemPriceCodeRespository.
func (r *ItemPriceCodeRepositoryImpl) GetAllItemPriceCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var results []masteritempayloads.SaveItemPriceCode
	var totalItems int64

	tableStruct := masteritementities.ItemPriceCode{}
	baseModelQuery := tx.Model(&tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	// Get the total number of items
	err := whereQuery.Count(&totalItems).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Use GetOffset and GetLimit for pagination
	err = whereQuery.Limit(pages.GetLimit()).Offset(pages.GetOffset()).Find(&results).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Convert results to map
	mapResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		mapResults[i] = map[string]interface{}{
			"price_code":      result.ItemPriceCodeId,
			"price_code_name": result.ItemPriceCodeName,
			"is_active":       result.IsActive,
			"item_price_id":   result.ItemPriceCode,
		}
	}

	return mapResults, int(totalItems), pages.Limit, nil
}

// GetByIdItemPriceCode retrieves an item price code by ID
func (r *ItemPriceCodeRepositoryImpl) GetByIdItemPriceCode(tx *gorm.DB, id int) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	var result masteritempayloads.SaveItemPriceCode
	var itemPriceCode masteritementities.ItemPriceCode

	// Query to find the ItemPriceCode by ID
	err := tx.Model(&itemPriceCode).Where("item_price_code_id = ?", id).First(&itemPriceCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item price code not found"),
			}
		}
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Map the entity to the response payload
	result.ItemPriceCodeId = itemPriceCode.ItemPriceCodeId
	result.IsActive = itemPriceCode.IsActive
	result.ItemPriceCode = itemPriceCode.ItemPriceCode
	result.ItemPriceCodeName = itemPriceCode.ItemPriceCodeName

	return result, nil
}

func (r *ItemPriceCodeRepositoryImpl) GetByCodeItemPriceCode(tx *gorm.DB, itemPriceCode string) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	var result masteritempayloads.SaveItemPriceCode
	var itemPriceCodeEntity masteritementities.ItemPriceCode

	// Query to find the ItemPriceCode by its code
	err := tx.Model(&itemPriceCodeEntity).Where("item_price_code = ?", itemPriceCode).First(&itemPriceCodeEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item price code not found"),
			}
		}
		return result, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Map the entity to the response payload
	result.ItemPriceCodeId = itemPriceCodeEntity.ItemPriceCodeId
	result.IsActive = itemPriceCodeEntity.IsActive
	result.ItemPriceCode = itemPriceCodeEntity.ItemPriceCode
	result.ItemPriceCodeName = itemPriceCodeEntity.ItemPriceCodeName

	return result, nil
}

// SaveItemPriceCode saves a new item price code to the database
func (r *ItemPriceCodeRepositoryImpl) SaveItemPriceCode(tx *gorm.DB, request masteritempayloads.SaveItemPriceCode) (masteritementities.ItemPriceCode, *exceptions.BaseErrorResponse) {
	newItemPriceCode := masteritementities.ItemPriceCode{
		ItemPriceCode:     request.ItemPriceCode,
		ItemPriceCodeName: request.ItemPriceCodeName,
		IsActive:          request.IsActive,
	}

	// Save the new ItemPriceCode
	err := tx.Model(&masteritementities.ItemPriceCode{}).Create(&newItemPriceCode).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return newItemPriceCode, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        errors.New("item price code already exists"),
			}
		}
		return newItemPriceCode, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return newItemPriceCode, nil
}

// DeleteItemPriceCode deletes an item price code by its ID
func (r *ItemPriceCodeRepositoryImpl) DeleteItemPriceCode(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	var itemPriceCode masteritementities.ItemPriceCode

	// Attempt to delete the ItemPriceCode
	err := tx.Model(&itemPriceCode).Where("item_price_code_id = ?", id).Delete(&itemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

// UpdateItemPriceCode updates an existing item price code
func (r *ItemPriceCodeRepositoryImpl) UpdateItemPriceCode(tx *gorm.DB, itemPriceId int, req masteritempayloads.UpdateItemPriceCode) (bool, *exceptions.BaseErrorResponse) {
	var existingItemPriceCode masteritementities.ItemPriceCode

	// Find the existing item price code by ID
	err := tx.Model(&existingItemPriceCode).Where("item_price_code_id = ?", itemPriceId).First(&existingItemPriceCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item price code not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Update the fields
	existingItemPriceCode.ItemPriceCode = req.ItemPriceCode
	existingItemPriceCode.ItemPriceCodeName = req.ItemPriceCodeName

	// Save the updated entity
	err = tx.Save(&existingItemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

// ChangeStatusItemPriceCode toggles the active status of an item price code
func (r *ItemPriceCodeRepositoryImpl) ChangeStatusItemPriceCode(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var existingItemPriceCode masteritementities.ItemPriceCode

	// Find the item price code by ID
	err := tx.Model(&existingItemPriceCode).Where("item_price_code_id = ?", id).First(&existingItemPriceCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("item price code not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Toggle the active status
	existingItemPriceCode.IsActive = !existingItemPriceCode.IsActive

	// Save the changes
	err = tx.Save(&existingItemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
