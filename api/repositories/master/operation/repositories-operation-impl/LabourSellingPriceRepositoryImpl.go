package masteroperationrepositoryimpl

import (
	"after-sales/api/config"
	masteroperationentities "after-sales/api/entities/master/operation"
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type LabourSellingPriceRepositoryImpl struct {
}

func StartLabourSellingPriceRepositoryImpl() masteroperationrepository.LabourSellingPriceRepository {
	return &LabourSellingPriceRepositoryImpl{}
}

// GetSellingPriceDetailById implements masteroperationrepository.LabourSellingPriceRepository.
func (r *LabourSellingPriceRepositoryImpl) GetSellingPriceDetailById(tx *gorm.DB, detailId int) (masteroperationpayloads.LabourSellingPriceDetailbyIdResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.LabourSellingPriceDetail{}
	response := masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}
	var modelResponse masteroperationpayloads.ModelSellingPriceDetailResponse
	var variantResponse masteroperationpayloads.VariantResponse

	if err := tx.Model(entities).Where(masteroperationentities.LabourSellingPriceDetail{LabourSellingPriceDetailId: detailId}).
		First(&entities).
		Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	response.LabourSellingPrice = entities.SellingPrice
	response.IsActive = entities.IsActive
	response.ModelId = entities.ModelId
	response.VariantId = entities.VariantId

	// join with mtr_unit_model

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entities.ModelId)

	errUrlUnitModel := utils.Get(unitModelUrl, &modelResponse, nil)

	if errUrlUnitModel != nil {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	if modelResponse == (masteroperationpayloads.ModelSellingPriceDetailResponse{}) {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New("model not found"),
		}
	}

	response.Model = modelResponse.ModelCode + " - " + modelResponse.ModelDescription

	//JOIN UNIT VARIANT

	unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entities.VariantId)

	errUrlunitVariant := utils.Get(unitVariantUrl, &variantResponse, nil)

	if errUrlunitVariant != nil {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	if variantResponse == (masteroperationpayloads.VariantResponse{}) {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New("varinat not found"),
		}
	}

	response.Variant = variantResponse.VariantCode + " - " + variantResponse.VariantDescription
	response.RecordStatus = variantResponse.VariantDescription

	return response, nil
}

// SaveMultipleDetail implements masteroperationrepository.LabourSellingPriceRepository.
func (r *LabourSellingPriceRepositoryImpl) SaveMultipleDetail(tx *gorm.DB, detail []masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptions.BaseErrorResponse) {

	for _, request := range detail {
		entities := masteroperationentities.LabourSellingPriceDetail{
			LabourSellingPriceId: request.LabourSellingPriceId,
			ModelId:              request.ModelId,
			VariantId:            request.VariantId,
			SellingPrice:         request.SellingPrice,
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
	}

	return true, nil

}

// GetAllDetailbyHeaderId implements masteroperationrepository.LabourSellingPriceRepository.
func (r *LabourSellingPriceRepositoryImpl) GetAllDetailbyHeaderId(tx *gorm.DB, headerId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := []masteroperationentities.LabourSellingPriceDetail{}
	responses := []masteroperationpayloads.LabourSellingPriceDetailResponse{}
	var getModelResponse []masteroperationpayloads.ModelSellingPriceDetailResponse
	var getVariantResponse []masteroperationpayloads.VariantResponse
	var ModelIds string
	var VariantIds string
	//define base model
	query := tx.
		Model(&entities).
		Where(masteroperationentities.LabourSellingPriceDetail{LabourSellingPriceId: headerId})

	//apply pagination and execute
	rows, err := query.Scan(&responses).Rows()

	if len(responses) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	models_ids := []int{}
	variant_ids := []int{}

	for _, response := range responses {
		if isNotInList(models_ids, response.ModelId) {
			str := strconv.Itoa(response.ModelId)
			ModelIds += str + ","
			models_ids = append(models_ids, response.ModelId)
		}
		if isNotInList(variant_ids, response.VariantId) {
			str := strconv.Itoa(response.VariantId)
			VariantIds += str + ","
			variant_ids = append(variant_ids, response.VariantId)
		}

	}

	// join with mtr_unit_model

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model-multi-id/" + ModelIds

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	if len(getModelResponse) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin(responses, getModelResponse, "ModelId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(getModelResponse) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	//JOIN UNIT VARIANT

	unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant-multi-id/" + VariantIds

	errUrlunitVariant := utils.Get(unitVariantUrl, &getVariantResponse, nil)

	if errUrlunitVariant != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	if len(getVariantResponse) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, getVariantResponse, "VariantId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(joinedData2) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	return joinedData2, nil
}

// GetAllSellingPrice implements masteroperationrepository.LabourSellingPriceRepository.
func (r *LabourSellingPriceRepositoryImpl) GetAllSellingPrice(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.LabourSellingPrice{}
	responses := []masteroperationentities.LabourSellingPrice{}

	query := tx.Model(entities)

	filterQuery := utils.ApplyFilterExact(query, filter)

	if err := filterQuery.Scopes(pagination.Paginate(entities, &pages, filterQuery)).Scan(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = responses

	return pages, nil

}

func isNotInList(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return false
		}
	}
	return true
}

func (r *LabourSellingPriceRepositoryImpl) GetLabourSellingPriceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.LabourSellingPrice{}
	response := masteroperationpayloads.LabourSellingPriceResponse{}
	var getUnitBrandResponse masteroperationpayloads.BrandLabourSellingPriceResponse
	var getjobTypeResponse masteroperationpayloads.JobTypeLabourSellingPriceResponse

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.LabourSellingPrice{
			LabourSellingPriceId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fmt.Print(response)

	defer rows.Close()

	// join with mtr_brand on sales service

	unitBrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(response.BrandId)

	errUrlUnitBrand := utils.Get(unitBrandUrl, &getUnitBrandResponse, nil)

	if errUrlUnitBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitBrand,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin([]masteroperationpayloads.LabourSellingPriceResponse{response}, []masteroperationpayloads.BrandLabourSellingPriceResponse{getUnitBrandResponse}, "BrandId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(joinedData1) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to fetch with brand"),
		}
	}

	//join with mtr_job_type on general service

	jobTypeUrl := config.EnvConfigs.GeneralServiceUrl + "job-type/" + strconv.Itoa(response.JobTypeId)

	errUrljobType := utils.Get(jobTypeUrl, &getjobTypeResponse, nil)

	if errUrljobType != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrljobType,
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, []masteroperationpayloads.JobTypeLabourSellingPriceResponse{getjobTypeResponse}, "JobTypeId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(joinedData2) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to fetch with job type"),
		}
	}

	result := joinedData2[0]

	return result, nil
}

func (r *LabourSellingPriceRepositoryImpl) GetAllSellingPriceDetailByHeaderId(tx *gorm.DB, headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entities := []masteroperationentities.LabourSellingPriceDetail{}
	responses := []masteroperationpayloads.LabourSellingPriceDetailResponse{}
	var getModelResponse []masteroperationpayloads.ModelSellingPriceDetailResponse
	var getVariantResponse []masteroperationpayloads.VariantResponse

	var ModelIds string
	var VariantIds string
	//define base model
	query := tx.
		Model(&entities).
		Where(masteroperationentities.LabourSellingPriceDetail{LabourSellingPriceId: headerId})

	fmt.Print(headerId)

	//apply pagination and execute
	rows, err := query.Scan(&responses).Rows()

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	models_ids := []int{}
	variant_ids := []int{}

	for _, response := range responses {
		if isNotInList(models_ids, response.ModelId) {
			str := strconv.Itoa(response.ModelId)
			ModelIds += str + ","
			models_ids = append(models_ids, response.ModelId)
		}
		if isNotInList(variant_ids, response.VariantId) {
			str := strconv.Itoa(response.VariantId)
			VariantIds += str + ","
			variant_ids = append(variant_ids, response.VariantId)
		}

	}

	// join with mtr_unit_model

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model-multi-id/" + ModelIds

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin(responses, getModelResponse, "ModelId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	//JOIN UNIT VARIANT

	unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant-multi-id/" + VariantIds

	errUrlunitVariant := utils.Get(unitVariantUrl, &getVariantResponse, nil)

	if errUrlunitVariant != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlUnitModel,
		}
	}

	if len(getVariantResponse) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, getVariantResponse, "VariantId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(joinedData2) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New(""),
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *LabourSellingPriceRepositoryImpl) SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (int, *exceptions.BaseErrorResponse) {

	entities := masteroperationentities.LabourSellingPrice{
		CompanyId:     request.CompanyId,
		BrandId:       request.BrandId,
		JobTypeId:     request.JobTypeId,
		EffectiveDate: request.EffectiveDate,
		BillToId:      request.BillToId,
		Description:   request.Description,
	}

	err := tx.Save(&entities).Where(entities).First(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return -1, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return -1, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities.LabourSellingPriceId, nil
}

func (r *LabourSellingPriceRepositoryImpl) SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (int, *exceptions.BaseErrorResponse) {

	entity_check := masteroperationentities.LabourSellingPriceDetail{}
	response := masteroperationpayloads.LabourSellingPriceDetailResponse{}
	LabourSellingPriceId := request.LabourSellingPriceId
	ModelId := request.ModelId
	VariantId := request.VariantId

	err1 := tx.Model(&entity_check).
		Where("labour_selling_price_id = ? AND model_id = ? AND variant_id = ?", LabourSellingPriceId, ModelId, VariantId).
		First(&response).
		Error

	if err1 == nil {
		return -1, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err1,
			Message: "Data already exist",
		}
	}


	entities := masteroperationentities.LabourSellingPriceDetail{
		LabourSellingPriceId: request.LabourSellingPriceId,
		ModelId:              request.ModelId,
		VariantId:            request.VariantId,
		SellingPrice:         request.SellingPrice,
	}

	err := tx.Save(&entities).Where(entities).First(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return -1, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return -1, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities.LabourSellingPriceDetailId, nil
}
