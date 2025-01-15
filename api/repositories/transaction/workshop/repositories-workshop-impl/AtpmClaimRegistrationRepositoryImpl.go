package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AtpmClaimRegistrationRepositoryImpl struct {
}

func OpenAtpmClaimRegistrationRepositoryImpl() transactionworkshoprepository.AtpmClaimRegistrationRepository {
	return &AtpmClaimRegistrationRepositoryImpl{}
}

// uspg_atAtpmVehicleClaim0_Select
// IF @Option = 0
func (r *AtpmClaimRegistrationRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

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

		vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
		if vehicleErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: vehicleErr.StatusCode,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        vehicleErr.Err,
			}
		}

		result := map[string]interface{}{
			"vehicle_id":                 entity.VehicleId,
			"vehicle_chassis_number":     vehicleResponse.Data.Master.VehicleChassisNumber,
			"work_order_system_number":   entity.WorkOrderSystemNumber,
			"work_order_document_number": entity.WorkOrderDocumentNumber,
			"work_order_date":            entity.WorkOrderDate,
			"company_id":                 entity.CompanyId,
			"company_name":               companyResponses.CompanyName,
			"claim_number":               entity.ClaimNumber,
			"claim_date":                 entity.ClaimDate,
			"claim_type_id":              entity.ClaimTypeId,
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

func (r *AtpmClaimRegistrationRepositoryImpl) GetById(tx *gorm.DB, id int, pages pagination.Pagination) (transactionworkshoppayloads.AtpmClaimRegistrationResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.AtpmClaimVehicle

	tx.Where("claim_system_number = ?", id).Find(&entity)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        tx.Error,
			}
		}
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        tx.Error,
		}
	}

	// Get company data
	companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
	if companyErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: companyErr.StatusCode,
			Message:    "Failed to fetch company data from internal service",
			Err:        companyErr.Err,
		}
	}

	// Get brand data
	brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: brandErr.StatusCode,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Get model data
	modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: modelErr.StatusCode,
			Message:    "Failed to fetch model data from external service",
			Err:        modelErr.Err,
		}
	}

	// Get variant data
	variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch variant data from external service",
			Err:        variantErr.Err,
		}
	}

	// Get vehicle data
	vehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.AtpmClaimRegistrationResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: vehicleErr.StatusCode,
			Message:    "Failed to fetch vehicle data from external service",
			Err:        vehicleErr.Err,
		}
	}

	response := transactionworkshoppayloads.AtpmClaimRegistrationResponse{
		ClaimSystemNumber:       entity.ClaimSystemNumber,
		CompanyId:               entity.CompanyId,
		CompanyName:             companyResponses.CompanyName,
		BrandId:                 entity.BrandId,
		BrandName:               brandResponses.BrandName,
		ClaimTypeId:             entity.ClaimTypeId,
		ClaimNumber:             entity.ClaimNumber,
		ClaimDate:               entity.ClaimDate,
		WorkOrderDocumentNumber: entity.WorkOrderDocumentNumber,
		WorkOrderDate:           entity.WorkOrderDate,
		VehicleId:               entity.VehicleId,
		VehicleChassisNumber:    vehicleResponse.Data.Master.VehicleChassisNumber,
		VehicleEngineNumber:     vehicleResponse.Data.Master.VehicleEngineNumber,
		ModelId:                 entity.ModelId,
		ModelDescription:        modelResponses.ModelName,
		VariantId:               entity.VariantId,
		VariantDescription:      variantResponses.VariantDescription,
	}

	return response, nil
}

func (r *AtpmClaimRegistrationRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.AtpmClaimRegistrationRequest) (transactionworkshopentities.AtpmClaimVehicle, *exceptions.BaseErrorResponse) {
	entity := transactionworkshopentities.AtpmClaimVehicle{
		CompanyId:            request.CompanyId,
		BrandId:              request.BrandId,
		ClaimTypeId:          request.ClaimTypeId,
		CustomerComplaint:    request.CustomerComplaint,
		TechnicianDiagnostic: request.TechnicianDiagnostic,
		Countermeasure:       request.Countermeasure,
		ClaimDate:            request.ClaimDate,
		RepairEndDate:        request.RepairEndDate,

		// other data
		Fuel:       request.Fuel,
		CustomerId: request.CustomerId,
		Vdn:        request.VDN,

		// Claim Header, Symptom, Trouble Code
		ClaimHeader: request.ClaimHeader,
	}

	if err := tx.Create(&entity).Error; err != nil {
		return transactionworkshopentities.AtpmClaimVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create data",
			Err:        tx.Error,
		}
	}

	return entity, nil
}
