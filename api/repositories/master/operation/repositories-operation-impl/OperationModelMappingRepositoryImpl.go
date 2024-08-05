package masteroperationrepositoryimpl

import (
	"after-sales/api/config"
	// masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
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

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationModelMappingResponse, *exceptions.BaseErrorResponse) {
	entities := masteroperationentities.OperationModelMapping{}
	response := masteroperationpayloads.OperationModelMappingResponse{}

	rows, err := tx.Model(&entities).
		Where(masteroperationentities.OperationModelMapping{
			OperationModelMappingId: Id,
		}).
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

func (r *OperationModelMappingRepositoryImpl) GetOperationModelMappingLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
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
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
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
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		operationModelMappingResponses = append(operationModelMappingResponses, response)
	}

	if len(operationModelMappingResponses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Fetch brand data
	var brandResponses []masteroperationpayloads.BrandResponse
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand?page=0&limit=10"
	errUrlBrand := utils.Get(brandUrl, &brandResponses, nil)
	if errUrlBrand != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlBrand,
		}
	}

	// Fetch model data
	var modelResponses []masteroperationpayloads.ModelResponse
	modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model?page=0&limit=10"
	errUrlModel := utils.Get(modelUrl, &modelResponses, nil)
	if errUrlModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
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
			return nil, 0, 0, &exceptions.BaseErrorResponse{
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

func (r *OperationModelMappingRepositoryImpl) SaveOperationModelMapping(tx *gorm.DB, request masteroperationpayloads.OperationModelMappingResponse) (bool, *exceptions.BaseErrorResponse) {
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

func (r *OperationModelMappingRepositoryImpl) GetAllOperationFrt(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = OperationFrtResponse

	return pages, nil
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

func (r *OperationModelMappingRepositoryImpl) GetAllOperationDocumentRequirement(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = OperationDocumentRequirementResponse

	return pages, nil
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
	OperationLevelMapping := []masteroperationentities.OperationLevel{}
	OperationLevelResponse := []masteroperationpayloads.OperationLevelRequest{}
	// OperationLevelResponse1 := masteroperationpayloads.OperationLevelResponse{}
	query := tx.
		Model(masteroperationentities.OperationLevel{}).
		Where("operation_model_mapping_id = ?", id).
		Scan(&OperationLevelResponse)

	err := query.
		Scopes(pagination.Paginate(&OperationLevelMapping, &pages, query)).
		Scan(&OperationLevelResponse).
		Error

	if len(OperationLevelResponse) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	// defer row.Close()
	pages.Rows = OperationLevelResponse

	return pages, nil
}

func (*OperationModelMappingRepositoryImpl) GetOperationLevelById(tx *gorm.DB, Id int) (masteroperationpayloads.OperationLevelByIdResponse, *exceptions.BaseErrorResponse) {
	var OperationLevelMapping masteroperationentities.OperationLevel
	var OperationLevelResponse masteroperationpayloads.OperationLevelByIdResponse

	rows, err := tx.
		Model(&OperationLevelMapping).
		Select(
			"operation_level_id",
			"OperationLevel.is_active is_active",
			"OperationEntries.operation_entries_code operation_entries_code",
			"OperationEntries.operation_entries_desc operation_entries_desc",
			"OperationGroup.operation_group_code operation_group_code",
			"OperationGroup.operation_group_description operation_group_description",
			"OperationSection.operation_section_code operation_section_code",
			"OperationSection.operation_section_description operation_section_description",
			"OperationKey.operation_key_code operation_key_code",
			"OperationKey.operation_key_description operation_key_description",
		).
		Joins("OperationEntries", tx.Select("1")).
		Joins("OperationEntries.OperationKey", tx.Select("1")).
		Joins("OperationEntries.OperationGroup", tx.Select("1")).
		Joins("OperationEntries.OperationSection", tx.Select("1")).
		Where(masteroperationentities.OperationLevel{OperationLevelId: Id}).
		First(&OperationLevelResponse).
		Rows()

	if err != nil {

		return OperationLevelResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	return OperationLevelResponse, nil
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

// func (r *OperationModelMappingRepositoryImpl) GetOperationLevel(tx *gorm.DB, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	// OperationLevelEtity := []masteroperationentities.OperationLevel{}
// 	OperationEntries := []masteroperationentities.OperationEntries{}
// 	OperationLevelPayloads := []masteroperationpayloads.OperationLevelGetAll{}

// 	_, err := tx.Model(&OperationEntries).
// 		Joins("Inner Join mtr_operation_group On mtr_operation_entries.operation_group_id=mtr_operation_group.operation_group_id").
// 		Joins("Inner Join mtr_operation_section On mtr_operation_section.operation_section_id=mtr_operation_entries.operation_section_id").
// 		Joins("Inner Join mtr_operation_key on mtr_operation_key.operation_key_id=mtr_operation_entries.operation_key_id").
// 		Group("mtr_operation_entries.operation_entries_code,mtr_operation_entries.operation_entries_description,mtr_operation_group.operation_group_code,mtr_operation_group.operation_group_description,mtr_operation_section.operation_section_code,mtr_operation_section.operation_section_description,mtr_operation_key.operation_key_code,mtr_operation_key.operation_key_code,mtr_operation_code_entries.is_active").
// 		Scan(&OperationLevelPayloads).Rows()
// 	if err != nil {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusNotFound,
// 			Err:        err,
// 		}
// 	}
// 	if len(OperationLevelPayloads) == 0 {
// 		return pages, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusNotFound,
// 			Err:        err,
// 		}
// 	}
// 	// defer query.close{}
// 	pages.Rows = OperationLevelPayloads
// 	return pages, nil
// }

// func (r *OperationModelMappingRepositoryImpl) DeleteOperationLevel(tx *gorm.DB, ids string)(bool,*exceptions.BaseErrorResponse){
// 	err:=tx.Model(masteroperationentities.OperationModelMapping).Delete()
// }
