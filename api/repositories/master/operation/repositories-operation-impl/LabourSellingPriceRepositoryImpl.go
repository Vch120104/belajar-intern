package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	"after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"

	"after-sales/api/utils"
	"errors"
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

	modelData, modelError := salesserviceapiutils.GetUnitModelById(entities.ModelId)
	if modelError != nil {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: modelError.StatusCode,
			Err:        modelError.Err,
		}
	}

	if modelData == (salesserviceapiutils.UnitModelResponse{}) {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New("model not found"),
		}
	}

	response.Model = modelData.ModelCode + " - " + modelData.ModelName

	variantData, variantError := salesserviceapiutils.GetUnitVariantById(entities.VariantId)
	if variantError != nil {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: variantError.StatusCode,
			Err:        variantError.Err,
		}
	}

	if variantData == (salesserviceapiutils.UnitVariantResponse{}) {
		return masteroperationpayloads.LabourSellingPriceDetailbyIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNoContent,
			Err:        errors.New("variant not found"),
		}
	}

	response.Variant = variantData.VariantCode + " - " + variantData.VariantDescription
	response.RecordStatus = variantData.VariantDescription

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

	query := tx.
		Model(&entities).
		Where(masteroperationentities.LabourSellingPriceDetail{LabourSellingPriceId: headerId})

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

	modelsIds := []int{}
	variantIds := []int{}
	ModelIds := ""
	VariantIds := ""

	for _, response := range responses {
		if isNotInList(modelsIds, response.ModelId) {
			ModelIds += strconv.Itoa(response.ModelId) + ","
			modelsIds = append(modelsIds, response.ModelId)
		}
		if isNotInList(variantIds, response.VariantId) {
			VariantIds += strconv.Itoa(response.VariantId) + ","
			variantIds = append(variantIds, response.VariantId)
		}
	}

	modelData, modelError := salesserviceapiutils.GetUnitModelByMultiId(modelsIds)
	if modelError != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: modelError.StatusCode,
			Err:        modelError.Err,
		}
	}

	variantData, variantError := salesserviceapiutils.GetUnitVariantByMultiId(variantIds)
	if variantError != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: variantError.StatusCode,
			Err:        variantError.Err,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin(responses, modelData, "ModelId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, variantData, "VariantId")
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
func (r *LabourSellingPriceRepositoryImpl) GetAllSellingPrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteroperationpayloads.LabourSellingPriceResponse

	baseModelQuery := tx.Model(&masteroperationentities.LabourSellingPrice{})
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
		// Fetch Brand data
		brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(response.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Err:        brandErr.Err,
			}
		}

		// Fetch Job Type data
		jobTypeResponse, jobTypeErr := generalserviceapiutils.GetJobTransactionTypeById(response.JobTypeId)
		if jobTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: jobTypeErr.StatusCode,
				Err:        jobTypeErr.Err,
			}
		}

		// Fetch BillTo (Supplier) data
		var billToName string
		if response.BillToId == 0 {
			billToName = ""
		} else {
			billToResponse, billToErr := generalserviceapiutils.GetSupplierMasterById(response.BillToId)
			if billToErr != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: billToErr.StatusCode,
					Err:        billToErr.Err,
				}
			}
			billToName = billToResponse.SupplierName
		}

		result := map[string]interface{}{
			"labour_selling_price_id": response.LabourSellingPriceId,
			"company_id":              response.CompanyId,
			"brand_id":                response.BrandId,
			"brand_name":              brandResponse.BrandName,
			"job_type_id":             response.JobTypeId,
			"job_type_name":           jobTypeResponse.JobTypeName,
			"effective_date":          response.EffectiveDate,
			"bill_to_id":              response.BillToId,
			"bill_to_name":            billToName,
			"description":             response.Description,
			"is_active":               response.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results

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
	var getUnitBrandResponse salesserviceapiutils.UnitBrandResponse
	var getJobTypeResponse generalserviceapiutils.WorkOrderJobType

	err := tx.Model(&entities).
		Where(masteroperationentities.LabourSellingPrice{
			LabourSellingPriceId: Id,
		}).
		First(&response).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	getUnitBrandResponse, errBrand := salesserviceapiutils.GetUnitBrandById(response.BrandId)
	if errBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: errBrand.StatusCode,
			Err:        errBrand.Err,
		}
	}

	getJobTypeResponse, errJobType := generalserviceapiutils.GetJobTransactionTypeById(response.JobTypeId)
	if errJobType != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: errJobType.StatusCode,
			Err:        errJobType.Err,
		}
	}

	joinedData1, errdf := utils.DataFrameInnerJoin([]masteroperationpayloads.LabourSellingPriceResponse{response}, []salesserviceapiutils.UnitBrandResponse{getUnitBrandResponse}, "BrandId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, []generalserviceapiutils.WorkOrderJobType{getJobTypeResponse}, "JobTypeId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	if len(joinedData2) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to fetch with brand and job type"),
		}
	}

	result := joinedData2[0]
	return result, nil
}

func (r *LabourSellingPriceRepositoryImpl) GetAllSellingPriceDetailByHeaderId(tx *gorm.DB, headerId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteroperationpayloads.LabourSellingPriceDetailResponse
	var results []map[string]interface{}

	modelIds := []int{}
	variantIds := []int{}

	query := tx.Model(&masteroperationentities.LabourSellingPriceDetail{}).
		Where(masteroperationentities.LabourSellingPriceDetail{LabourSellingPriceId: headerId})

	err := query.Scopes(pagination.Paginate(&pages, query)).Find(&responses).Error
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

	for _, response := range responses {
		modelIds = append(modelIds, response.ModelId)
		variantIds = append(variantIds, response.VariantId)
	}

	modelData, modelErr := salesserviceapiutils.GetUnitModelByMultiId(modelIds)
	if modelErr != nil {
		return pages, modelErr
	}

	joinedData1, joinErr := utils.DataFrameInnerJoin(responses, modelData, "ModelId")
	if joinErr != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        joinErr,
		}
	}

	if len(variantIds) > 0 {
		variantData, variantErr := salesserviceapiutils.GetUnitVariantByMultiId(variantIds)
		if variantErr != nil {
			return pages, variantErr
		}

		joinedData2, joinErr := utils.DataFrameInnerJoin(joinedData1, variantData, "VariantId")
		if joinErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        joinErr,
			}
		}

		for _, data := range joinedData2 {
			results = append(results, map[string]interface{}{
				"labour_selling_price_id":        data["LabourSellingPriceId"],
				"model_id":                       data["ModelId"],
				"model_code":                     data["ModelCode"],
				"model_description":              data["ModelDescription"],
				"variant_id":                     data["VariantId"],
				"variant_code":                   data["VariantCode"],
				"variant_description":            data["VariantDescription"],
				"selling_price":                  data["SellingPrice"],
				"labour_selling_price_detail_id": data["LabourSellingPriceDetailId"],
				"is_active":                      data["IsActive"],
			})
		}
	} else {
		for _, data := range joinedData1 {
			results = append(results, map[string]interface{}{
				"labour_selling_price_id":        data["LabourSellingPriceId"],
				"model_id":                       data["ModelId"],
				"model_code":                     data["ModelCode"],
				"model_description":              data["ModelDescription"],
				"selling_price":                  data["SellingPrice"],
				"labour_selling_price_detail_id": data["LabourSellingPriceDetailId"],
				"is_active":                      data["IsActive"],
			})
		}
	}

	pages.Rows = results
	return pages, nil
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
			Message:    "Data already exist",
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

func (r *LabourSellingPriceRepositoryImpl) DeleteLabourSellingPriceDetail(tx *gorm.DB, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	var entities []masteroperationentities.LabourSellingPriceDetail

	result := tx.Where("labour_selling_price_detail_id IN ?", iddet).Find(&entities)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if err := tx.Delete(&entities).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}
