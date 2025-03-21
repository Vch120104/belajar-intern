package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"net/http"

	"gorm.io/gorm"
)

type WarrantyFreeServiceRepositoryImpl struct {
}

func StartWarrantyFreeServiceRepositoryImpl() masterrepository.WarrantyFreeServiceRepository {
	return &WarrantyFreeServiceRepositoryImpl{}
}

func (r *WarrantyFreeServiceRepositoryImpl) GetAllWarrantyFreeService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masterpayloads.WarrantyFreeServiceResponse

	baseModelQuery := tx.Model(&masterentities.WarrantyFreeService{})
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
				StatusCode: http.StatusInternalServerError,
				Err:        brandErr.Err,
			}
		}

		// Fetch Model data
		modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(response.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        modelErr.Err,
			}
		}

		// Fetch Warranty Free Service Type data
		warrantyFreeServiceTypeResponse, wfstErr := generalserviceapiutils.GetWarrantyFreeServiceTypeById(response.WarrantyFreeServiceTypeId)
		if wfstErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        wfstErr.Err,
			}
		}

		result := map[string]interface{}{
			"warranty_free_service_id":               response.WarrantyFreeServicesId,
			"brand_id":                               response.BrandId,
			"brand_name":                             brandResponse.BrandName,
			"model_id":                               response.ModelId,
			"model_name":                             modelResponse.ModelName,
			"warranty_free_service_type_id":          response.WarrantyFreeServiceTypeId,
			"warranty_free_service_type_code":        warrantyFreeServiceTypeResponse.WarrantyFreeServiceTypeCode,
			"warranty_free_service_type_description": warrantyFreeServiceTypeResponse.WarrantyFreeServiceTypeName,
			"effective_date":                         response.EffectiveDate,
			"expire_mileage":                         response.ExpireMileage,
			"expire_month":                           response.ExpireMonth,
			"variant_id":                             response.VariantId,
			"expire_mileage_extended_warranty":       response.ExpireMileageExtendedWarranty,
			"expire_month_extended_warranty":         response.ExpireMonthExtendedWarranty,
			"remark":                                 response.Remark,
			"extended_warranty":                      response.ExtendedWarranty,
			"is_active":                              response.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) GetWarrantyFreeServiceById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := masterentities.WarrantyFreeService{}
	response := masterpayloads.WarrantyFreeServiceResponse{}

	// Fetch data from database
	err := tx.Model(&entities).
		Where("warranty_free_services_id = ?", Id).
		First(&response).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch Brand data
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(response.BrandId)
	if brandErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        brandErr.Err,
		}
	}

	// Fetch Model data
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(response.ModelId)
	if modelErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        modelErr.Err,
		}
	}

	// Fetch Variant data
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(response.VariantId)
	if variantErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        variantErr.Err,
		}
	}

	// Fetch Warranty Free Service Type data
	warrantyFreeServiceTypeResponse, wfstErr := generalserviceapiutils.GetWarrantyFreeServiceTypeById(response.WarrantyFreeServiceTypeId)
	if wfstErr != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        wfstErr.Err,
		}
	}

	// Prepare the response as a map
	finalResponse := map[string]interface{}{
		"is_active":                        response.IsActive,
		"warranty_free_services_id":        response.WarrantyFreeServicesId,
		"brand_id":                         response.BrandId,
		"model_id":                         response.ModelId,
		"warranty_free_service_type_id":    response.WarrantyFreeServiceTypeId,
		"effective_date":                   response.EffectiveDate,
		"expire_mileage":                   response.ExpireMileage,
		"expire_month":                     response.ExpireMonth,
		"variant_id":                       response.VariantId,
		"expire_mileage_extended_warranty": response.ExpireMileageExtendedWarranty,
		"expire_month_extended_warranty":   response.ExpireMonthExtendedWarranty,
		"remark":                           response.Remark,
		"extended_warranty":                response.ExtendedWarranty,
		"brand_code":                       brandResponse.BrandCode,
		"model_code":                       modelResponse.ModelCode,
		"variant_code":                     variantResponse.VariantCode,
		"brand_name":                       brandResponse.BrandName,
		"model_description":                modelResponse.ModelName,
		"variant_description":              variantResponse.VariantDescription,
		"warranty_free_service_type_code":  warrantyFreeServiceTypeResponse.WarrantyFreeServiceTypeCode,
		"warranty_free_service_type_name":  warrantyFreeServiceTypeResponse.WarrantyFreeServiceTypeName,
	}

	return finalResponse, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) SaveWarrantyFreeService(tx *gorm.DB, request masterpayloads.WarrantyFreeServiceRequest) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse) {
	entities := masterentities.WarrantyFreeService{
		BrandId:                       request.BrandId,
		ModelId:                       request.ModelId,
		WarrantyFreeServiceTypeId:     request.WarrantyFreeServiceTypeId,
		EffectiveDate:                 request.EffectiveDate,
		ExpireMileage:                 &request.ExpireMileage,
		ExpireMonth:                   &request.ExpireMonth,
		VariantId:                     request.VariantId,
		ExpireMileageExtendedWarranty: &request.ExpireMileageExtendedWarranty,
		ExpireMonthExtendedWarranty:   &request.ExpireMonthExtendedWarranty,
		Remark:                        &request.Remark,
		ExtendedWarranty:              &request.ExtendedWarranty,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return masterentities.WarrantyFreeService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) ChangeStatusWarrantyFreeService(tx *gorm.DB, Id int) (masterpayloads.WarrantyFreeServicePatchResponse, *exceptions.BaseErrorResponse) {
	var entities masterentities.WarrantyFreeService
	var response masterpayloads.WarrantyFreeServicePatchResponse

	result := tx.Model(&entities).
		Where(masterentities.WarrantyFreeService{WarrantyFreeServicesId: Id}).
		First(&entities)

	if result.Error != nil {
		return response, &exceptions.BaseErrorResponse{
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
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	data := tx.Model(&entities).
		Where(masterentities.WarrantyFreeService{WarrantyFreeServicesId: Id}).
		First(&response)

	if data.Error != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return response, nil
}

func (r *WarrantyFreeServiceRepositoryImpl) UpdateWarrantyFreeService(tx *gorm.DB, req masterentities.WarrantyFreeService, id int) (masterentities.WarrantyFreeService, *exceptions.BaseErrorResponse) {
	var entity masterentities.WarrantyFreeService

	err := tx.Model(entity).Where(masterentities.WarrantyFreeService{WarrantyFreeServicesId: id}).Updates(req).First(&entity).Error
	if err != nil {
		return masterentities.WarrantyFreeService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity, nil
}
