package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AtpmReimbursementRepositoryImpl struct {
}

func OpenAtpmReimbursementRepositoryImpl() *AtpmReimbursementRepositoryImpl {
	return &AtpmReimbursementRepositoryImpl{}
}

// uspg_atAtpmVehicleClaim0_Select
// IF @Option = 0
func (r *AtpmReimbursementRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshopentities.AtpmClaimVehicle

	tx = utils.ApplyFilter(tx.Model(&transactionworkshopentities.AtpmClaimVehicle{}), filterCondition)

	tx.Scopes(pagination.Paginate(&pages, tx)).Find(&entities)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        tx.Error,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        tx.Error,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		// vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
		// if vehicleErr != nil {
		// 	return pages, &exceptions.BaseErrorResponse{
		// 		StatusCode: vehicleErr.StatusCode,
		// 		Message:    "Failed to fetch vehicle data from external service",
		// 		Err:        vehicleErr.Err,
		// 	}
		// }

		// fetch claim type
		claimTypeResponse, claimTypeErr := generalserviceapiutils.GetClaimTypeById(entity.ClaimTypeId)
		if claimTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: claimTypeErr.StatusCode,
				Message:    "Failed to fetch claim type data from internal service",
				Err:        claimTypeErr.Err,
			}
		}

		result := map[string]interface{}{
			"vehicle_id":                 entity.VehicleId,
			"vehicle_chassis_number":     "vehicleResponse.Data.Master.VehicleChassisNumber",
			"work_order_system_number":   entity.WorkOrderSystemNumber,
			"work_order_document_number": entity.WorkOrderDocumentNumber,
			"work_order_date":            entity.WorkOrderDate,
			"company_id":                 entity.CompanyId,
			"company_name":               companyResponses.CompanyName,
			"claim_system_number":        entity.ClaimSystemNumber,
			"claim_number":               entity.ClaimNumber,
			"claim_date":                 entity.ClaimDate,
			"claim_type_id":              entity.ClaimTypeId,
			"claim_type_description":     claimTypeResponse.ClaimTypeDescription,
			"claim_status_id":            entity.ClaimStatusId,
			"brand_id":                   entity.BrandId,
			"brand_name":                 brandResponses.BrandName,
			"model_id":                   entity.ModelId,
			"model_description":          modelResponses.ModelName,
			"variant_id":                 entity.VariantId,
			"variant_description":        variantResponses.VariantDescription,
			"claim_from":                 entity.ClaimFrom,
			"claim_to":                   entity.ClaimTo,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}
