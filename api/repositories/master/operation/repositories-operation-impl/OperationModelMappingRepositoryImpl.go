package masteroperationrepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	aftersalesserviceapiutils "after-sales/api/utils/aftersales-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"

	exceptions "after-sales/api/exceptions"
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

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingById(tx *gorm.DB, id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
	var response masteroperationpayloads.OperationModelMappingResponse

	err := tx.Table("mtr_operation_model_mapping").
		Select("mtr_operation_model_mapping.is_active, "+
			"mtr_operation_model_mapping.operation_model_mapping_id, "+
			"mtr_operation_model_mapping.brand_id, "+
			"mtr_operation_model_mapping.model_id, "+
			"mtr_operation_model_mapping.operation_id, "+
			"mtr_operation_code.operation_code, "+
			"mtr_operation_code.operation_name, "+
			"mtr_operation_model_mapping.operation_using_incentive AS operation_using_incentive, "+
			"mtr_operation_model_mapping.operation_using_actual AS operation_using_actual, "+
			"mtr_operation_model_mapping.operation_pdi AS operation_pdi").
		Joins("JOIN mtr_operation_code ON mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id").
		Where("mtr_operation_model_mapping.operation_model_mapping_id = ?", id).
		First(&response).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get data",
			Err:        err,
		}
	}

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingByBrandModelOperationCode(tx *gorm.DB, request masteroperationpayloads.OperationModelModelBrandOperationCodeRequest) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
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
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []masteroperationentities.OperationModelMapping
	tx = tx.
		Select("mtr_operation_model_mapping.operation_model_mapping_id, " +
			"mtr_operation_model_mapping.brand_id, mtr_operation_model_mapping.model_id, " +
			"mtr_operation_model_mapping.operation_id, mtr_operation_model_mapping.is_active, " +
			"mtr_operation_code.operation_code, mtr_operation_code.operation_name").
		Joins("JOIN mtr_operation_code ON mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id")
	tx = utils.ApplyFilter(tx.Model(&masteroperationentities.OperationModelMapping{}), filterCondition)
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
		brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		operationResponse, operationErr := aftersalesserviceapiutils.GetOperationById(entity.OperationId)
		if operationErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: operationErr.StatusCode,
				Message:    "Failed to fetch operation data from external service",
				Err:        operationErr.Err,
			}
		}

		result := map[string]interface{}{
			"operation_model_mapping_id": entity.OperationModelMappingId,
			"brand_id":                   entity.BrandId,
			"brand_name":                 brandResponse.BrandName,
			"model_id":                   entity.ModelId,
			"model_code":                 modelResponse.ModelCode,
			"operation_id":               entity.OperationId,
			"operation_code":             operationResponse.OperationCode,
			"operation_name":             operationResponse.OperationName,
			"is_active":                  entity.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMapping(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationModelMapping{
		IsActive:                request.IsActive,
		OperationModelMappingId: request.OperationModelMappingId,
		BrandId:                 request.BrandId,
		ModelId:                 request.ModelId,
		OperationId:             request.OperationId,
		OperationUsingIncentive: &request.OperationUsingIncentive,
		OperationUsingActual:    &request.OperationUsingActual,
		OperationPdi:            &request.OperationPdi,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ChangeStatusOperationModelMapping(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationModelMapping

	result := tx.Model(&entities).
		Where("operation_model_mapping_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationFrt(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var OperationFrtResponse []masteroperationpayloads.OperationModelMappingFrtRequest
	var variantIds []int

	err := tx.Table(masteroperationentities.TableNameOperationFrt).
		Select("mtr_operation_frt.operation_frt_id AS operation_frt_id, "+
			"mtr_operation_frt.operation_model_mapping_id AS operation_model_mapping_id, "+
			"mtr_operation_frt.variant_id AS variant_id, "+
			"mtr_operation_frt.frt_hour AS frt_hour, "+
			"mtr_operation_frt.frt_hour_express AS frt_hour_express, "+
			"mtr_operation_frt.is_active AS is_active").
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationFrtResponse).Error

	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, op := range OperationFrtResponse {
		variantIds = append(variantIds, op.VariantId)
	}

	variantData, errVariant := salesserviceapiutils.GetUnitVariantByMultiId(variantIds)
	if errVariant != nil {
		return nil, 0, 0, errVariant
	}

	joinedData, errJoin := utils.DataFrameInnerJoin(OperationFrtResponse, variantData, "variant_id")
	if errJoin != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errJoin,
		}
	}

	results, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return results, totalPages, totalRows, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationFrtByHeaderId(tx *gorm.DB, id int) ([]masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse) {
	var OperationFrtResponse []masteroperationpayloads.OperationModelMappingFrtRequest

	err := tx.Table(masteroperationentities.TableNameOperationFrt).
		Select("mtr_operation_frt.operation_frt_id AS operation_frt_id, "+
			"mtr_operation_frt.operation_model_mapping_id AS operation_model_mapping_id, "+
			"mtr_operation_frt.variant_id AS variant_id, "+
			"mtr_operation_frt.frt_hour AS frt_hour, "+
			"mtr_operation_frt.frt_hour_express AS frt_hour_express, "+
			"mtr_operation_frt.is_active AS is_active").
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationFrtResponse).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return OperationFrtResponse, nil
}

func (*OperationModelMappingRepositoryImpl) GetOperationFrtById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingFrtRequest, *exceptions.BaseErrorResponse) {
	var OperationFrtMapping masteroperationentities.OperationFrt
	var OperationFrtResponse masteroperationpayloads.OperationModelMappingFrtRequest

	rows, err := tx.
		Model(&OperationFrtMapping).
		Where(masteroperationentities.OperationFrt{OperationFrtId: Id}).
		First(&OperationFrtResponse).
		Rows()

	if err != nil {

		return OperationFrtResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	return OperationFrtResponse, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMappingFrt(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingFrtRequest) (bool, *exceptions.BaseErrorResponse) {
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) DeleteOperationLevel(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteroperationentities.OperationLevel
	if err := tx.Delete(&entities, ids).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) DeactivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationFrt
		err := tx.Model(&entityToUpdate).Where("operation_frt_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ActivateOperationFrt(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationFrt
		err := tx.Model(&entityToUpdate).Where("operation_frt_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationDocumentRequirement(
	tx *gorm.DB,
	id int,
	pages pagination.Pagination,
) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var OperationDocumentRequirementResponse []masteroperationpayloads.OperationModelMappingDocumentRequirementRequest

	// Start query on the table
	query := tx.
		Model(&masteroperationentities.OperationDocumentRequirement{}).
		Where("operation_model_mapping_id = ?", id)

	// Apply pagination to the query using the Paginate function
	queryWithPagination := query.Scopes(pagination.Paginate(&pages, query))

	// Execute the query
	err := queryWithPagination.Scan(&OperationDocumentRequirementResponse).Error

	// Handle the error if any
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If no records found
	if len(OperationDocumentRequirementResponse) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Assign the paginated rows to the pages.Rows field
	pages.Rows = OperationDocumentRequirementResponse

	return pages, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationDocumentRequirementByHeaderId(tx *gorm.DB, id int) ([]masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse) {

	var OperationDocumentRequirementResponse []masteroperationpayloads.OperationModelMappingDocumentRequirementRequest

	err := tx.
		Model(masteroperationentities.OperationDocumentRequirement{}).
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationDocumentRequirementResponse).Error

	if err != nil {

		return OperationDocumentRequirementResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return OperationDocumentRequirementResponse, nil
}

func (*OperationModelMappingRepositoryImpl) GetOperationDocumentRequirementById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingDocumentRequirementRequest, *exceptions.BaseErrorResponse) {
	var OperationDocumentRequirementMapping masteroperationentities.OperationDocumentRequirement
	var OperationDocumentRequirementResponse masteroperationpayloads.OperationModelMappingDocumentRequirementRequest

	rows, err := tx.
		Model(&OperationDocumentRequirementMapping).
		Where(masteroperationentities.OperationDocumentRequirement{OperationDocumentRequirementId: Id}).
		First(&OperationDocumentRequirementResponse).
		Rows()

	if err != nil {

		return OperationDocumentRequirementResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	return OperationDocumentRequirementResponse, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMappingDocumentRequirement(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingDocumentRequirementRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationDocumentRequirement{
		IsActive:                                request.IsActive,
		OperationModelMappingId:                 request.OperationModelMappingId,
		OperationDocumentRequirementId:          request.OperationDocumentRequirementId,
		Line:                                    request.Line,
		OperationDocumentRequirementDescription: request.OperationDocumentRequirementDescription,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) DeactivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationDocumentRequirement
		err := tx.Model(&entityToUpdate).Where("operation_document_requirement_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ActivateOperationDocumentRequirement(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationDocumentRequirement
		err := tx.Model(&entityToUpdate).Where("operation_document_requirement_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) SaveOperationLevel(tx *gorm.DB, request masteroperationpayloads.OperationLevelRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationLevel{
		IsActive:                request.IsActive,
		OperationLevelId:        request.OperationLevelId,
		OperationModelMappingId: request.OperationModelMappingId,
		OperationEntriesId:      request.OperationEntriesId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationLevel(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var OperationLevelResponse []masteroperationpayloads.OperationLevelGetAll

	query := tx.Table("mtr_operation_level").Select(`
		mtr_operation_level.operation_level_id,
		mtr_operation_group.operation_group_id,
		mtr_operation_group.operation_group_code,
		mtr_operation_group.operation_group_description,
		mtr_operation_section.operation_section_id,
		mtr_operation_section.operation_section_code,
		mtr_operation_section.operation_section_description,
		mtr_operation_key.operation_key_id,
		mtr_operation_key.operation_key_code,
		mtr_operation_key.operation_key_description,
		op_entries.operation_entries_id,
		op_entries.operation_entries_code,
		op_entries.operation_entries_description`).
		Joins("JOIN mtr_operation_entries AS op_entries ON op_entries.operation_entries_id = mtr_operation_level.operation_entries_id").
		Joins("JOIN mtr_operation_group ON mtr_operation_group.operation_group_id = op_entries.operation_group_id").
		Joins("JOIN mtr_operation_key ON mtr_operation_key.operation_key_id = op_entries.operation_key_id").
		Joins("JOIN mtr_operation_section ON mtr_operation_section.operation_section_id = op_entries.operation_section_id").
		Where("mtr_operation_level.operation_model_mapping_id = ?", id)

	err := query.Scopes(pagination.Paginate(&pages, query)).Scan(&OperationLevelResponse).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If no data is found, return empty array
	if len(OperationLevelResponse) == 0 {
		pages.Rows = []masteroperationpayloads.OperationLevelGetAll{}
		return pages, nil
	}

	pages.Rows = OperationLevelResponse

	return pages, nil
}

func (r *OperationModelMappingRepositoryImpl) GetAllOperationLevelByHeaderId(tx *gorm.DB, id int) ([]masteroperationpayloads.OperationLevelRequest, *exceptions.BaseErrorResponse) {

	var OperationLevelResponse []masteroperationpayloads.OperationLevelRequest

	err := tx.Table("mtr_operation_level").Select(`
		mtr_operation_level.operation_level_id,
		mtr_operation_group.operation_group_id,
		mtr_operation_group.operation_group_code,
		mtr_operation_group.operation_group_description,
		mtr_operation_section.operation_section_id,
		mtr_operation_section.operation_section_code,
		mtr_operation_section.operation_section_description,
		mtr_operation_key.operation_key_id,
		mtr_operation_key.operation_key_code,
		mtr_operation_key.operation_key_description,
		op_entries.operation_entries_id,
		op_entries.operation_entries_code,
		op_entries.operation_entries_description`).
		Joins("JOIN mtr_operation_entries AS op_entries ON op_entries.operation_entries_id = mtr_operation_level.operation_entries_id").
		Joins("JOIN mtr_operation_group ON mtr_operation_group.operation_group_id = op_entries.operation_group_id").
		Joins("JOIN mtr_operation_key ON mtr_operation_key.operation_key_id = op_entries.operation_key_id").
		Joins("JOIN mtr_operation_section ON mtr_operation_section.operation_section_id = op_entries.operation_section_id").
		Where("mtr_operation_level.operation_model_mapping_id = ?", id).
		Scan(&OperationLevelResponse).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return OperationLevelResponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Level not found",
				Err:        err,
			}
		}
		return OperationLevelResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving Operation Level",
			Err:        err,
		}
	}

	return OperationLevelResponse, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationLevelById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationLevelByIdResponse, *exceptions.BaseErrorResponse) {
	response := masteroperationpayloads.OperationLevelByIdResponse{}

	err := tx.Model(&masteroperationentities.OperationLevel{}).
		Select(`
			mtr_operation_level.operation_level_id,
			mtr_operation_level.is_active,
			mtr_operation_model_mapping.operation_model_mapping_id,
			mtr_operation_entries.operation_entries_id,
			mtr_operation_entries.operation_entries_code,
			mtr_operation_entries.operation_entries_description,
			mtr_operation_group.operation_group_code,
			mtr_operation_group.operation_group_description,
			mtr_operation_section.operation_section_code,
			mtr_operation_section.operation_section_description,
			mtr_operation_key.operation_key_code,
			mtr_operation_key.operation_key_description`).
		Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_operation_level.operation_model_mapping_id").
		Joins("JOIN mtr_operation_entries ON mtr_operation_entries.operation_entries_id = mtr_operation_level.operation_entries_id").
		Joins("JOIN mtr_operation_key ON mtr_operation_key.operation_key_id = mtr_operation_entries.operation_key_id").
		Joins("JOIN mtr_operation_group ON mtr_operation_group.operation_group_id = mtr_operation_entries.operation_group_id").
		Joins("JOIN mtr_operation_section ON mtr_operation_section.operation_section_id = mtr_operation_entries.operation_section_id").
		Where("mtr_operation_level.operation_level_id = ?", Id).
		First(&response).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Level not found",
				Err:        err,
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving Operation Level",
			Err:        err,
		}
	}

	return response, nil
}

func (r *OperationModelMappingRepositoryImpl) DeactivateOperationLevel(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationLevel
		err := tx.Model(&entityToUpdate).Where("operation_level_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = false
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) ActivateOperationLevel(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	idSlice := strings.Split(id, ",")

	for _, Ids := range idSlice {
		var entityToUpdate masteroperationentities.OperationLevel
		err := tx.Model(&entityToUpdate).Where("operation_level_id = ?", Ids).First(&entityToUpdate).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		entityToUpdate.IsActive = true
		result := tx.Save(&entityToUpdate)
		if result.Error != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
	}

	return true, nil
}

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingLatestId(tx *gorm.DB) (int, *exceptions.BaseErrorResponse) {
	var latestID int

	err := tx.Table("mtr_operation_model_mapping").
		Select("operation_model_mapping_id").
		Order("operation_model_mapping_id DESC").
		Limit(1).
		Scan(&latestID).Error

	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return latestID, nil
}

func (r *OperationModelMappingRepositoryImpl) UpdateOperationModelMapping(tx *gorm.DB, operationModelMappingId int, request masteroperationpayloads.OperationModelMappingUpdate) (masteroperationentities.OperationModelMapping, *exceptions.BaseErrorResponse) {
	var OperationModelMapping = masteroperationentities.OperationModelMapping{
		OperationPdi:            &request.OperationPdi,
		OperationUsingIncentive: &request.OperationUsingIncentive,
		OperationUsingActual:    &request.OperationUsingActual,
	}

	if err := tx.Model(&masteroperationentities.OperationModelMapping{}).
		Where("operation_model_mapping_id = ?", operationModelMappingId).
		Updates(&OperationModelMapping).Error; err != nil {
		return masteroperationentities.OperationModelMapping{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update operation master",
			Err:        err,
		}
	}

	return OperationModelMapping, nil
}

func (r *OperationModelMappingRepositoryImpl) UpdateOperationFrt(tx *gorm.DB, operationFrtId int, request masteroperationpayloads.OperationFrtUpdate) (masteroperationentities.OperationFrt, *exceptions.BaseErrorResponse) {
	var OperationFrt = masteroperationentities.OperationFrt{
		OperationFrtId: request.OperationFrtId,
		FrtHour:        request.FrtHour,
		FrtHourExpress: request.FrtHourExpress,
	}

	// if err := tx.Model(&masteroperationentities.OperationFrt{}).
	// 	Where("operation_frt_id = ?", operationFrtId).
	// 	Updates(&OperationFrt).Error; err != nil {
	// 	return masteroperationentities.OperationFrt{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to update operation master",
	// 		Err:        err,
	// 	}
	// }

	if err := tx.Model(&masteroperationentities.OperationFrt{}).
		Where("operation_frt_id = ?", operationFrtId).
		Select("FrtHour", "FrtHourExpress").
		Updates(&OperationFrt).Error; err != nil {
		return masteroperationentities.OperationFrt{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update operation master",
			Err:        err,
		}
	}

	return OperationFrt, nil
}

func (r *OperationModelMappingRepositoryImpl) CopyOperationModelMappingToOtherModel(tx *gorm.DB, headerId int, request masteroperationpayloads.OperationModelMappingCopyRequest) (bool, *exceptions.BaseErrorResponse) {

	var CopyRequest masteroperationpayloads.CopyRequest
	var latestID int

	operationModelMappingHeader, _ := r.GetOperationModelMappingById(tx, headerId)

	CopyRequest.HeaderRequest = operationModelMappingHeader

	operationLevel, _ := r.GetAllOperationLevelByHeaderId(tx, headerId)

	CopyRequest.OperationLevel = operationLevel

	operationFrt, _ := r.GetAllOperationFrtByHeaderId(tx, headerId)

	CopyRequest.OperationFrt = operationFrt

	operationDoc, _ := r.GetAllOperationDocumentRequirementByHeaderId(tx, headerId)

	CopyRequest.OperationDoc = operationDoc

	headerEntity := masteroperationentities.OperationModelMapping{
		IsActive:                CopyRequest.HeaderRequest.IsActive,
		OperationModelMappingId: 0,
		BrandId:                 CopyRequest.HeaderRequest.BrandId,
		ModelId:                 request.ModelId,
		OperationId:             CopyRequest.HeaderRequest.OperationId,
		OperationUsingIncentive: &CopyRequest.HeaderRequest.OperationUsingIncentive,
		OperationUsingActual:    &CopyRequest.HeaderRequest.OperationUsingActual,
		OperationPdi:            &CopyRequest.HeaderRequest.OperationPdi,
	}

	err := tx.Save(&headerEntity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	errLatestId := tx.Table("mtr_operation_model_mapping").
		Select("operation_model_mapping_id").
		Order("operation_model_mapping_id DESC").
		Limit(1).
		Scan(&latestID).Error

	if errLatestId != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errLatestId,
		}
	}

	CopyRequest.HeaderRequest.OperationModelMappingId = latestID

	for _, operationFrtValues := range CopyRequest.OperationFrt {
		oprFrtEntity := masteroperationentities.OperationFrt{
			IsActive:                operationFrtValues.IsActive,
			OperationFrtId:          0,
			OperationModelMappingId: latestID,
			VariantId:               operationFrtValues.VariantId,
			FrtHour:                 operationFrtValues.FrtHour,
			FrtHourExpress:          operationFrtValues.FrtHourExpress,
		}

		err := tx.Create(&oprFrtEntity).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	for _, operationLevelValues := range CopyRequest.OperationLevel {
		oprLvlEntity := masteroperationentities.OperationLevel{
			IsActive:                operationLevelValues.IsActive,
			OperationLevelId:        0,
			OperationModelMappingId: latestID,
			OperationEntriesId:      operationLevelValues.OperationEntriesId,
		}

		err := tx.Create(&oprLvlEntity).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	for _, operationDocValues := range CopyRequest.OperationDoc {
		oprDocEntity := masteroperationentities.OperationDocumentRequirement{
			IsActive:                                operationDocValues.IsActive,
			OperationModelMappingId:                 latestID,
			OperationDocumentRequirementId:          0,
			Line:                                    operationDocValues.Line,
			OperationDocumentRequirementDescription: operationDocValues.OperationDocumentRequirementDescription,
		}

		err := tx.Create(&oprDocEntity).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}
