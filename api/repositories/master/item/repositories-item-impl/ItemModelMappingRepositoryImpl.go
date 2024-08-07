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
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type ItemModelMappingRepositoryImpl struct {
}

// GetItemModelMappingByItemId implements masteritemrepository.ItemModelMappingRepository.
func (r *ItemModelMappingRepositoryImpl) GetItemModelMappingByItemId(tx *gorm.DB, itemId int, pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemModelMappingReponses
	var brandResponses []masteritempayloads.UnitBrandResponses
	var modelResponses []masteritempayloads.UnitModelResponses
	var variantResponses []masteritempayloads.UnitVariantResponses

	model := masteritementities.ItemDetail{}
	baseModelQuery := tx.Model(&model).Where(masteritementities.ItemDetail{ItemId: itemId})

	rows, err := baseModelQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	//// Brand
	brandUrl := config.EnvConfigs.SalesServiceUrl + "/unit-brand-dropdown"

	if errBrand := utils.Get(brandUrl, &brandResponses, nil); errBrand != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedDataBrand := utils.DataFrameInnerJoin(responses, brandResponses, "BrandId")

	if len(joinedDataBrand) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	//// Unit Model
	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "/unit-model-dropdown"

	if errModel := utils.Get(unitModelUrl, &modelResponses, nil); errModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedDataModel := utils.DataFrameInnerJoin(joinedDataBrand, modelResponses, "ModelId")

	if len(joinedDataModel) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	//// Unit Variant
	unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "/unit-variant?page=0&limit=1000"

	if errVariant := utils.Get(unitVariantUrl, &variantResponses, nil); errVariant != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	joinedDataVariant := utils.DataFrameInnerJoin(joinedDataModel, variantResponses, "VariantId")

	if len(joinedDataVariant) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedDataVariant, &pages)

	return dataPaginate, totalPages, totalRows, nil

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
