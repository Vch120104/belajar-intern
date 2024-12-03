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

func (r *ItemPriceCodeRepositoryImpl) GetAllItemPriceCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.SaveItemPriceCode

	baseModelQuery := tx.Model(&masteritementities.ItemPriceCode{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {
		result := map[string]interface{}{
			"item_price_code_id":   response.ItemPriceCodeId,
			"item_price_code_name": response.ItemPriceCodeName,
			"is_active":            response.IsActive,
			"item_price_code":      response.ItemPriceCode,
		}
		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *ItemPriceCodeRepositoryImpl) GetByIdItemPriceCode(tx *gorm.DB, id int) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	var result masteritempayloads.SaveItemPriceCode
	var itemPriceCode masteritementities.ItemPriceCode

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

	result.ItemPriceCodeId = itemPriceCode.ItemPriceCodeId
	result.IsActive = itemPriceCode.IsActive
	result.ItemPriceCode = itemPriceCode.ItemPriceCode
	result.ItemPriceCodeName = itemPriceCode.ItemPriceCodeName

	return result, nil
}

func (r *ItemPriceCodeRepositoryImpl) GetByCodeItemPriceCode(tx *gorm.DB, itemPriceCode string) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	var result masteritempayloads.SaveItemPriceCode
	var itemPriceCodeEntity masteritementities.ItemPriceCode

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

	result.ItemPriceCodeId = itemPriceCodeEntity.ItemPriceCodeId
	result.IsActive = itemPriceCodeEntity.IsActive
	result.ItemPriceCode = itemPriceCodeEntity.ItemPriceCode
	result.ItemPriceCodeName = itemPriceCodeEntity.ItemPriceCodeName

	return result, nil
}

func (r *ItemPriceCodeRepositoryImpl) SaveItemPriceCode(tx *gorm.DB, request masteritempayloads.SaveItemPriceCode) (masteritementities.ItemPriceCode, *exceptions.BaseErrorResponse) {
	newItemPriceCode := masteritementities.ItemPriceCode{
		ItemPriceCode:     request.ItemPriceCode,
		ItemPriceCodeName: request.ItemPriceCodeName,
		IsActive:          request.IsActive,
	}

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

func (r *ItemPriceCodeRepositoryImpl) DeleteItemPriceCode(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var itemPriceCode masteritementities.ItemPriceCode

	err := tx.Model(&itemPriceCode).Where("item_price_code_id = ?", id).Delete(&itemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemPriceCodeRepositoryImpl) UpdateItemPriceCode(tx *gorm.DB, itemPriceId int, req masteritempayloads.UpdateItemPriceCode) (bool, *exceptions.BaseErrorResponse) {
	var existingItemPriceCode masteritementities.ItemPriceCode

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

	existingItemPriceCode.ItemPriceCode = req.ItemPriceCode
	existingItemPriceCode.ItemPriceCodeName = req.ItemPriceCodeName

	err = tx.Save(&existingItemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemPriceCodeRepositoryImpl) ChangeStatusItemPriceCode(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var existingItemPriceCode masteritementities.ItemPriceCode

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

	existingItemPriceCode.IsActive = !existingItemPriceCode.IsActive

	err = tx.Save(&existingItemPriceCode).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemPriceCodeRepositoryImpl) GetItemPriceCodeDropDown(tx *gorm.DB) ([]masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse) {
	model := masteritementities.ItemPriceCode{}
	responses := []masteritempayloads.SaveItemPriceCode{}
	err := tx.Model(model).Scan(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return responses, nil
}
