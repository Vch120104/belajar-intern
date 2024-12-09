package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemModelMappingRepositoryImpl struct {
}

// GetItemModelMappingByItemId implements masteritemrepository.ItemModelMappingRepository.
func (r *ItemModelMappingRepositoryImpl) GetItemModelMappingByItemId(tx *gorm.DB, itemId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemModelMappingReponses

	baseModelQuery := tx.Model(&masteritementities.ItemDetail{}).Where("item_id = ?", itemId)
	whereQuery := utils.ApplyFilter(baseModelQuery, nil)

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
		// Fetch Brand data
		brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(response.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data",
				Err:        brandErr.Err,
			}
		}

		// Fetch Model data
		modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data",
				Err:        modelErr.Err,
			}
		}

		// Fetch Variant data
		variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(response.VariantId)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: variantErr.StatusCode,
				Message:    "Failed to fetch variant data",
				Err:        variantErr.Err,
			}
		}

		result := map[string]interface{}{
			"item_id":      response.ItemId,
			"brand_id":     response.BrandId,
			"brand_name":   brandResponse.BrandName,
			"model_id":     response.ModelId,
			"model_name":   modelResponse.ModelName,
			"variant_id":   response.VariantId,
			"variant_name": variantResponse.VariantName,
			"is_active":    response.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

// UpdateItemModelMapping implements masteritemrepository.ItemModelMappingRepository.
func (r *ItemModelMappingRepositoryImpl) UpdateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {

	entities := masteritementities.ItemDetail{
		ItemDetailId: req.ItemDetailId,
		MileageEvery: req.MileageEvery,
		ReturnEvery:  req.ReturnEvery,
	}

	err := tx.Updates(&entities).Error

	if err != nil {

		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

// CreateItemModelMapping implements masteritemrepository.ItemModelMappingRepository.
func (r *ItemModelMappingRepositoryImpl) CreateItemModelMapping(tx *gorm.DB, req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {

	_, brandErr := salesserviceapiutils.GetUnitBrandById(req.BrandId)
	if brandErr != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: brandErr.StatusCode,
			Message:    "Failed to fetch brand data",
			Err:        brandErr.Err,
		}
	}

	_, modelErr := salesserviceapiutils.GetUnitModelById(req.ModelId)
	if modelErr != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: modelErr.StatusCode,
			Message:    "Failed to fetch model data",
			Err:        modelErr.Err,
		}
	}

	_, variantErr := salesserviceapiutils.GetUnitVariantById(req.VariantId)
	if variantErr != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: variantErr.StatusCode,
			Message:    "Failed to fetch variant data",
			Err:        variantErr.Err,
		}
	}

	entities := masteritementities.ItemDetail{
		IsActive:     req.IsActive,
		ItemId:       req.ItemId,
		BrandId:      req.BrandId,
		ModelId:      req.ModelId,
		VariantId:    req.VariantId,
		MileageEvery: req.MileageEvery,
		ReturnEvery:  req.ReturnEvery,
	}

	err := tx.Create(&entities).Error

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

func StartItemModelMappingRepositoryImpl() masteritemrepository.ItemModelMappingRepository {
	return &ItemModelMappingRepositoryImpl{}
}
