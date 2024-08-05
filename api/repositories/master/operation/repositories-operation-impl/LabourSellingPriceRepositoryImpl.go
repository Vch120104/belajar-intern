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

// GetAllSellingPrice implements masteroperationrepository.LabourSellingPriceRepository.
func (r *LabourSellingPriceRepositoryImpl) GetAllSellingPrice(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.LabourSellingPrice{}
	responses := []masteroperationentities.LabourSellingPrice{}

	query := tx.Model(entities)

	filterQuery := utils.ApplyFilter(query, filter)

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

	fmt.Println(unitBrandUrl)

	joinedData1 := utils.DataFrameInnerJoin([]masteroperationpayloads.LabourSellingPriceResponse{response}, []masteroperationpayloads.BrandLabourSellingPriceResponse{getUnitBrandResponse}, "BrandId")

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

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []masteroperationpayloads.JobTypeLabourSellingPriceResponse{getjobTypeResponse}, "JobTypeId")

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

	joinedData1 := utils.DataFrameInnerJoin(responses, getModelResponse, "ModelId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData1, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *LabourSellingPriceRepositoryImpl) SaveLabourSellingPrice(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptions.BaseErrorResponse) {

	entities := masteroperationentities.LabourSellingPrice{
		CompanyId:     request.CompanyId,
		BrandId:       request.BrandId,
		JobTypeId:     request.JobTypeId,
		EffectiveDate: request.EffectiveDate,
		BillToId:      request.BillToId,
		Description:   request.Description,
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

func (r *LabourSellingPriceRepositoryImpl) SaveLabourSellingPriceDetail(tx *gorm.DB, request masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptions.BaseErrorResponse) {

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

	return true, nil
}
