package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type WarrantyFreeServiceRepositoryImpl struct {
}

func StartWarrantyFreeServiceRepositoryImpl() masterrepository.WarrantyFreeServiceRepository {
	return &WarrantyFreeServiceRepositoryImpl{}
}

func (r *WarrantyFreeServiceRepositoryImpl) GetAllWarrantyFreeService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masterpayloads.WarrantyFreeServiceListResponse
	var getBrandResponse []masterpayloads.BrandResponse
	var getModelResponse []masterpayloads.UnitModelResponse
	var getWarrantyFreeServiceTypeResponse []masterpayloads.WarrantyFreeServiceTypeResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var brandCode string
	var modelCode string
	var warrantyFreeServiceTypeCode string
	responseStruct := reflect.TypeOf(masterpayloads.WarrantyFreeServiceListResponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
		}
	}

	//apply external services filter
	for i := 0; i < len(externalServiceFilter); i++ {
		if strings.Contains(externalServiceFilter[i].ColumnField, "brand_code") {
			brandCode = externalServiceFilter[i].ColumnValue
		} else if strings.Contains(externalServiceFilter[i].ColumnField, "warranty_free_service_type_code") {
			warrantyFreeServiceTypeCode = externalServiceFilter[i].ColumnValue
		} else {
			modelCode = externalServiceFilter[i].ColumnValue
		}
	}

	result := tx.Model(masterentities.WarrantyFreeService{})

	// define table struct
	// tableStruct := masterpayloads.WarrantyFreeServiceListResponse{}
	//define join table
	// joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(result, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// join with mtr_brand

	unitBrandUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-brand?page=0&limit=1000000&brand_code=" + brandCode

	errUrlUnitBrand := utils.Get(unitBrandUrl, &getBrandResponse, nil)

	if errUrlUnitBrand != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin(responses, getBrandResponse, "BrandId")

	// join with mtr_unit_model

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-model?page=0&limit=100000&model_code=" + modelCode

	errUrlUnitModel := utils.Get(unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, getModelResponse, "ModelId")

	// join with mtr_warranty_free_service_type

	warrantyFreeServiceTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/warranty-free-service-type?warranty_free_service_type_code=" + warrantyFreeServiceTypeCode

	errUrlWarrantyFreeServiceType := utils.Get(warrantyFreeServiceTypeUrl, &getWarrantyFreeServiceTypeResponse, nil)

	if errUrlWarrantyFreeServiceType != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, getWarrantyFreeServiceTypeResponse, "WarrantyFreeServiceTypeId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData3, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.WarrantyFreeService{}
	response := masterpayloads.WarrantyFreeServiceResponse{}
	var getUnitBrandResponse masterpayloads.BrandResponse
	var getUnitModelResponse masterpayloads.UnitModelResponse
	var getUnitVariantResponse masterpayloads.UnitVariantResponse
	var getWarrantyFreeServiceTypeResponse masterpayloads.WarrantyFreeServiceTypeResponse

	rows, err := tx.Model(&entities).
		Where(masterentities.WarrantyFreeService{
			WarrantyFreeServicesId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	// join with mtr_brand on sales service

	unitBrandUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-brand/" + strconv.Itoa(response.BrandId)

	errUrlUnitBrand := utils.Get(unitBrandUrl, &getUnitBrandResponse, nil)

	if errUrlUnitBrand != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData1 := utils.DataFrameInnerJoin([]masterpayloads.WarrantyFreeServiceResponse{response}, []masterpayloads.BrandResponse{getUnitBrandResponse}, "BrandId")

	//join with mtr_unit_model on sales service

	unitModelUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-model/" + strconv.Itoa(response.ModelId)

	errUrlUnitModel := utils.Get(unitModelUrl, &getUnitModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.UnitModelResponse{getUnitModelResponse}, "ModelId")

	// join with mtr_unit_variant on sales service

	unitVariantUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-variant/" + strconv.Itoa(response.VariantId)

	errUrlUnitVariant := utils.Get(unitVariantUrl, &getUnitVariantResponse, nil)

	if errUrlUnitVariant != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, []masterpayloads.UnitVariantResponse{getUnitVariantResponse}, "VariantId")

	// join with mtr_warranty_free_service_type on general service

	warrantyFreeServiceTypeUrl := config.EnvConfigs.GeneralServiceUrl + "/api/general/warranty-free-service-type/" + strconv.Itoa(response.WarrantyFreeServiceTypeId)

	errUrlWarrantyFreeServiceType := utils.Get(warrantyFreeServiceTypeUrl, &getWarrantyFreeServiceTypeResponse, nil)

	if errUrlWarrantyFreeServiceType != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData4 := utils.DataFrameInnerJoin(joinedData3, []masterpayloads.WarrantyFreeServiceTypeResponse{getWarrantyFreeServiceTypeResponse}, "WarrantyFreeServiceTypeId")

	result := joinedData4[0]

	return result, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.WarrantyFreeService{
		WarrantyFreeServicesId:        request.WarrantyFreeServicesId,
		BrandId:                       request.BrandId,
		ModelId:                       request.ModelId,
		WarrantyFreeServiceTypeId:     request.WarrantyFreeServiceTypeId,
		EffectiveDate:                 request.EffectiveDate,
		ExpireMileage:                 request.ExpireMileage,
		ExpireMonth:                   request.ExpireMonth,
		VariantId:                     request.VariantId,
		ExpireMileageExtendedWarranty: request.ExpireMileageExtendedWarranty,
		ExpireMonthExtendedWarranty:   request.ExpireMonthExtendedWarranty,
		Remark:                        request.Remark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) ChangeStatusWarrantyFreeService(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.WarrantyFreeService

	result := tx.Model(&entities).
		Where("warranty_free_services_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
