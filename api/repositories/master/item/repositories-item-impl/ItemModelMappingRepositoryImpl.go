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

	var wg sync.WaitGroup

	//// Brand
	wg.Add(1)
	go func() *exceptions.BaseErrorResponse {

		defer wg.Done()

		brandUrl := config.EnvConfigs.SalesServiceUrl + "/unit-brand-dropdown"

		if errBrand := utils.Get(brandUrl, &brandResponses, nil); errBrand != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(""),
			}
		}
		return nil
	}()

	//// Unit Model
	wg.Add(1)
	go func() *exceptions.BaseErrorResponse {

		defer wg.Done()

		unitModelUrl := config.EnvConfigs.SalesServiceUrl + "/unit-model-dropdown"

		if errModel := utils.Get(unitModelUrl, &modelResponses, nil); errModel != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(""),
			}
		}
		return nil
	}()

	//// Unit Variant
	wg.Add(1)
	go func() *exceptions.BaseErrorResponse {

		defer wg.Done()

		ids := ""

		for _, value := range responses {
			ids += fmt.Sprintf("%d,", value.VariantId)
		}

		unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant-multi-id/" + ids

		if errVariant := utils.Get(unitVariantUrl, &variantResponses, nil); errVariant != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(""),
			}
		}
		return nil
	}()

	//AWAIT Goroutines
	wg.Wait()

	joinedDataBrand, errdf := utils.DataFrameInnerJoin(responses, brandResponses, "BrandId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	joinedDataModel, errdf := utils.DataFrameInnerJoin(joinedDataBrand, modelResponses, "ModelId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	joinedDataVariant, errdf := utils.DataFrameInnerJoin(joinedDataModel, variantResponses, "VariantId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
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

	//Check brandID
	var brandResponses masteritempayloads.UnitBrandResponses
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(req.BrandId)

	if errBrandUrl := utils.Get(brandUrl, &brandResponses, nil); errBrandUrl != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("brand not found"),
		}
	}
	//

	//check unit model
	var unitmodelresponses masteritempayloads.UnitModelResponses

	unitmodelurl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(req.ModelId)

	if errunitmodelurl := utils.Get(unitmodelurl, &unitmodelresponses, nil); errunitmodelurl != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("unit model not found"),
		}
	}

	//check variant
	var variantResponses masteritempayloads.UnitVariantResponses

	variantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(req.VariantId)

	if errvarianturl := utils.Get(variantUrl, &variantResponses, nil); errvarianturl != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("variant not found"),
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
