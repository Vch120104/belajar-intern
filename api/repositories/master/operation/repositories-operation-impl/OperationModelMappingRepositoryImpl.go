package masteroperationrepositoryimpl

import (
	"after-sales/api/config"
	masteroperationentities "after-sales/api/entities/master/operation"
	exceptionsss_test "after-sales/api/expectionsss"
	"errors"
	"net/http"
	"strings"
	"time"

	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	"after-sales/api/utils"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OperationModelMappingRepositoryImpl struct {
}

func StartOperationModelMappingRepositoryImpl() masteroperationrepository.OperationModelMappingRepository {
	return &OperationModelMappingRepositoryImpl{}
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationModelMapping{}
	response := masteroperationpayloads.OperationModelMappingResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationModelMapping{
			OperationModelMappingId: Id,
		}).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingByBrandModelOperationCode(tx *gorm.DB, request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationModelMapping{}
	response := masteroperationpayloads.OperationModelMappingResponse{}

	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	tx.Logger = newLogger

	rows, err := tx.Model(&entities).
		Where("brand_id = ? AND model_id = ? AND operation_id = ?", request.BrandId, request.ModelId, request.OperationId).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []map[string]interface{}

	// Fetch OperationModelMapping data
	var operationModelMappingResponses []masteroperationpayloads.OperationModelMappingLookup

	// Define table struct
	tableStruct := masteroperationpayloads.OperationModelMappingLookup{}

	// Join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filter
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	// Fetch and map the data to OperationModelMappingLookup struct
	for rows.Next() {
		var response masteroperationpayloads.OperationModelMappingLookup
		if err := rows.Scan(
			&response.IsActive,
			&response.OperationModelMappingId,
			&response.OperationId,
			&response.OperationName,
			&response.OperationCode,
			&response.BrandId,
			&response.ModelId,
		); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		operationModelMappingResponses = append(operationModelMappingResponses, response)
	}

	if len(operationModelMappingResponses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Fetch brand data
	var brandResponses []masteroperationpayloads.BrandResponse
	brandUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-brand?page=0&limit=10"
	errUrlBrand := utils.Get(brandUrl, &brandResponses, nil)
	if errUrlBrand != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrand,
		}
	}

	// Fetch model data
	var modelResponses []masteroperationpayloads.ModelResponse
	modelUrl := config.EnvConfigs.SalesServiceUrl + "/api/sales/unit-model?page=0&limit=10"
	errUrlModel := utils.Get(modelUrl, &modelResponses, nil)
	if errUrlModel != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModel,
		}
	}

	// Create a map to hold brand and model data
	brandMap := make(map[int]masteroperationpayloads.BrandResponse)
	modelMap := make(map[int]masteroperationpayloads.ModelResponse)

	// Fill brand and model maps
	for _, brand := range brandResponses {
		brandMap[brand.BrandId] = brand
	}

	for _, model := range modelResponses {
		modelMap[model.ModelId] = model
	}

	// Combine data from OperationModelMapping, Brand, and Model
	for _, mapping := range operationModelMappingResponses {
		brand, brandExists := brandMap[mapping.BrandId]
		model, modelExists := modelMap[mapping.ModelId]

		if !brandExists || !modelExists {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("brand or model data not found"),
			}
		}

		response := map[string]interface{}{
			"IsActive":                mapping.IsActive,
			"OperationModelMappingId": mapping.OperationModelMappingId,
			"OperationId":             mapping.OperationId,
			"OperationName":           mapping.OperationName,
			"OperationCode":           mapping.OperationCode,
			"BrandId":                 mapping.BrandId,
			"ModelId":                 mapping.ModelId,
			"BrandName":               brand.BrandName,
			"ModelDescription":        model.ModelDescription,
		}

		responses = append(responses, response)
	}

	// Paginate the data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMapping(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationModelMapping{
		IsActive:                request.IsActive,
		OperationModelMappingId: request.OperationModelMappingId,
		BrandId:                 request.BrandId,
		ModelId:                 request.ModelId,
		OperationId:             request.OperationId,
		OperationUsingIncentive: request.OperationUsingIncentive,
		OperationUsingActual:    request.OperationUsingActual,
		OperationPdi:            request.OperationPdi,
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

func (r *OperationModelMappingRepositoryImpl) ChangeStatusOperationModelMapping(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masteroperationentities.OperationModelMapping

	result := tx.Model(&entities).
		Where("operation_model_mapping_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	// Toggle the IsActive value
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

func (r *OperationModelMappingRepositoryImpl) GetAllOperationFrt(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	OperationFrtMapping := []masteroperationentities.OperationFrt{}
	OperationFrtResponse := []masteroperationpayloads.OperationModelMappingFrtRequest{}

	query := tx.
		Model(masteroperationentities.OperationFrt{}).
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationFrtResponse)

	err := query.
		Scopes(pagination.Paginate(&OperationFrtMapping, &pages, query)).
		Scan(&OperationFrtResponse).
		Error

	if len(OperationFrtResponse) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = OperationFrtResponse

	return pages, nil
}

func (*OperationModelMappingRepositoryImpl) GetOperationFrtById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptionsss_test.BaseErrorResponse) {
	var OperationFrtMapping masteroperationentities.OperationFrt
	var OperationFrtResponse masteroperationpayloads.OperationModelMappingFrtRequest

	rows, err := tx.
		Model(&OperationFrtMapping).
		Where(masteroperationentities.OperationFrt{OperationFrtId: Id}).
		First(&OperationFrtResponse).
		Rows()

	if err != nil {

		return OperationFrtResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return OperationFrtResponse, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMappingFrt(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationFrt{
		IsActive:                request.IsActive,
		OperationFrtId:          request.OperationFrtId,
		OperationModelMappingId: request.OperationModelMappingId,
		VariantId:               request.VariantId,
		FrtHour:                 request.FrtHour,
		FrtHourExpress:          request.FrtHourExpress,
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

func (r *OperationModelMappingRepositoryImpl) DeactivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationFrt
		err := tx.Model(&entityToUpdate).Where("operation_frt_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ActivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationFrt
		err := tx.Model(&entityToUpdate).Where("operation_frt_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationDocumentRequirement(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	OperationDocumentRequirementMapping := []masteroperationentities.OperationDocumentRequirement{}
	OperationDocumentRequirementResponse := []masteroperationpayloads.OperationModelMappingDocumentRequirementRequest{}
	// OperationDocumentRequirementResponse1 := masteroperationpayloads.OperationDocumentRequirementResponse{}
	query := tx.
		Model(masteroperationentities.OperationDocumentRequirement{}).
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationDocumentRequirementResponse)

	err := query.
		Scopes(pagination.Paginate(&OperationDocumentRequirementMapping, &pages, query)).
		// Order("approval.name").
		Scan(&OperationDocumentRequirementResponse).
		Error

	if len(OperationDocumentRequirementResponse) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = OperationDocumentRequirementResponse

	return pages, nil
}

func (*OperationModelMappingRepositoryImpl) GetOperationDocumentRequirementById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptionsss_test.BaseErrorResponse) {
	var OperationDocumentRequirementMapping masteroperationentities.OperationDocumentRequirement
	var OperationDocumentRequirementResponse masteroperationpayloads.OperationModelMappingDocumentRequirementRequest

	rows, err := tx.
		Model(&OperationDocumentRequirementMapping).
		Where(masteroperationentities.OperationDocumentRequirement{OperationDocumentRequirementId: Id}).
		First(&OperationDocumentRequirementResponse).
		Rows()

	if err != nil {

		return OperationDocumentRequirementResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return OperationDocumentRequirementResponse, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMappingDocumentRequirement(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteroperationentities.OperationDocumentRequirement{
		IsActive:                                request.IsActive,
		OperationModelMappingId:                 request.OperationModelMappingId,
		OperationDocumentRequirementId:          request.OperationDocumentRequirementId,
		Line:                                    request.Line,
		OperationDocumentRequirementDescription: request.OperationDocumentRequirementDescription,
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

func (r *OperationModelMappingRepositoryImpl) DeactivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationDocumentRequirement
		err := tx.Model(&entityToUpdate).Where("operation_document_requirement_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ActivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptionsss_test.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationDocumentRequirement
		err := tx.Model(&entityToUpdate).Where("operation_document_requirement_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}
