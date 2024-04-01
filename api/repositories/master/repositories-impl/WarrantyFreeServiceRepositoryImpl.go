package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarrantyFreeServiceRepositoryImpl struct {
}

func StartWarrantyFreeServiceRepositoryImpl() masterrepository.WarrantyFreeServiceRepository {
	return &WarrantyFreeServiceRepositoryImpl{}
}

func (r *WarrantyFreeServiceRepositoryImpl) GetAllWarrantyFreeService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masterpayloads.WarrantyFreeServiceListResponse
	var getBrandResponse []masterpayloads.BrandResponse
	var getModelResponse []masterpayloads.UnitModelResponse
	var getWarrantyFreeServiceTypeResponse []masterpayloads.WarrantyFreeServiceTypeResponse
	var c *gin.Context
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
		} else if strings.Contains(externalServiceFilter[i].ColumnField, "warranty_free_service_type_code"){
			warrantyFreeServiceTypeCode = externalServiceFilter[i].ColumnValue
		}else {
			modelCode = externalServiceFilter[i].ColumnValue
		}
	}

	// define table struct
	tableStruct := masterpayloads.WarrantyFreeServiceListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, gorm.ErrRecordNotFound
	}

	// join with mtr_brand

	unitBrandUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand?page=0&limit=1000000&brand_code=" + brandCode

	errUrlUnitBrand := utils.Get(c, unitBrandUrl, &getBrandResponse, nil)

	if errUrlUnitBrand != nil {
		return nil, 0, 0, errUrlUnitBrand
	}

	joinedData1 := utils.DataFrameInnerJoin(responses, getBrandResponse, "BrandId")

	// join with mtr_unit_model

	unitModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model?page=0&limit=100000&model_code=" + modelCode

	errUrlUnitModel := utils.Get(c, unitModelUrl, &getModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, 0, 0, errUrlUnitModel
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, getModelResponse, "ModelId")

	// join with mtr_warranty_free_service_type

	warrantyFreeServiceTypeUrl := "http://10.1.32.26:8000/general-service/api/general/warranty-free-service-type?warranty_free_service_type_code=" + warrantyFreeServiceTypeCode

	errUrlWarrantyFreeServiceType := utils.Get(c, warrantyFreeServiceTypeUrl, &getWarrantyFreeServiceTypeResponse, nil)

	if errUrlWarrantyFreeServiceType != nil {
		return nil, 0, 0, errUrlWarrantyFreeServiceType
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, getWarrantyFreeServiceTypeResponse, "WarrantyFreeServiceTypeId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData3, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, error) {
	entities := masterentities.WarrantyFreeService{}
	response := masterpayloads.WarrantyFreeServiceResponse{}
	var c *gin.Context
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
		return nil, err
	}

	defer rows.Close()

	// join with mtr_brand on sales service

	unitBrandUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/" + strconv.Itoa(response.BrandId)

	errUrlUnitBrand := utils.Get(c, unitBrandUrl, &getUnitBrandResponse, nil)

	if errUrlUnitBrand != nil {
		return nil, errUrlUnitBrand
	}

	joinedData1 := utils.DataFrameInnerJoin([]masterpayloads.WarrantyFreeServiceResponse{response}, []masterpayloads.BrandResponse{getUnitBrandResponse}, "BrandId")

	//join with mtr_unit_model on sales service

	unitModelUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-model/" + strconv.Itoa(response.ModelId)

	errUrlUnitModel := utils.Get(c, unitModelUrl, &getUnitModelResponse, nil)

	if errUrlUnitModel != nil {
		return nil, errUrlUnitModel
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.UnitModelResponse{getUnitModelResponse}, "ModelId")

	// join with mtr_unit_variant on sales service

	unitVariantUrl := "http://10.1.32.26:8000/sales-service/api/sales/unit-variant/" + strconv.Itoa(response.VariantId)

	errUrlUnitVariant := utils.Get(c, unitVariantUrl, &getUnitVariantResponse, nil)

	if errUrlUnitVariant != nil {
		return nil, errUrlUnitVariant
	}

	joinedData3 := utils.DataFrameInnerJoin(joinedData2, []masterpayloads.UnitVariantResponse{getUnitVariantResponse}, "VariantId")

	// join with mtr_warranty_free_service_type on general service

	warrantyFreeServiceTypeUrl := "http://10.1.32.26:8000/general-service/api/general/warranty-free-service-type/" + strconv.Itoa(response.WarrantyFreeServiceTypeId)

	errUrlWarrantyFreeServiceType := utils.Get(c, warrantyFreeServiceTypeUrl, &getWarrantyFreeServiceTypeResponse, nil)

	if errUrlWarrantyFreeServiceType != nil {
		return nil, errUrlWarrantyFreeServiceType
	}

	joinedData4 := utils.DataFrameInnerJoin(joinedData3, []masterpayloads.WarrantyFreeServiceTypeResponse{getWarrantyFreeServiceTypeResponse}, "WarrantyFreeServiceTypeId")

	result := joinedData4[0]

	return result, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (bool, error) {
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
		return false, err
	}

	return true, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) ChangeStatusWarrantyFreeService(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.WarrantyFreeService

	result := tx.Model(&entities).
		Where("warranty_free_services_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	if entities.IsActive {
		entities.IsActive = false
	} else {
		entities.IsActive = true
	}

	result = tx.Save(&entities)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}