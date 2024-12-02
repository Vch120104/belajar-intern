package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type IncentiveMasterRepositoryImpl struct {
}

func StartIncentiveMasterRepositoryImpl() masterrepository.IncentiveMasterRepository {
	return &IncentiveMasterRepositoryImpl{}
}

func (r *IncentiveMasterRepositoryImpl) GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []masterentities.IncentiveMaster

	// Apply filters and pagination
	baseModelQuery := tx.Model(&masterentities.IncentiveMaster{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	// Perform the query with pagination
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// If no entities are found, return empty rows
	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		// Fetch the role data using the Role API
		role, errResp := generalserviceapiutils.GetRoleById(entity.JobPositionId)
		if errResp != nil {
			return pages, errResp
		}

		// Prepare the result map with the incentive details and role data
		result := map[string]interface{}{
			"incentive_level_id":      entity.IncentiveLevelId,
			"incentive_level_code":    entity.IncentiveLevelCode,
			"job_position_id":         entity.JobPositionId,
			"job_position_name":       role.RoleName, // Using RoleName from the RoleResponse
			"incentive_level_percent": entity.IncentiveLevelPercent,
			"is_active":               entity.IsActive,
		}

		results = append(results, result)
	}

	// Attach the results to the pagination rows
	pages.Rows = results
	return pages, nil
}

func (r *IncentiveMasterRepositoryImpl) GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, *exceptions.BaseErrorResponse) {
	entities := masterentities.IncentiveMaster{}
	response := masterpayloads.IncentiveMasterResponse{}

	err := tx.Model(&entities).
		Where(masterentities.IncentiveMaster{
			IncentiveLevelId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *IncentiveMasterRepositoryImpl) SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.IncentiveMaster{
		IncentiveLevelId:      request.IncentiveLevelId,
		IncentiveLevelCode:    request.IncentiveLevelCode,
		JobPositionId:         request.JobPositionId,
		IncentiveLevelPercent: request.IncentiveLevelPercent,
		IsActive:              request.IsActive,
	}

	if request.IncentiveLevelId == 0 {
		// Jika IncentiveMasterId == 0, ini adalah operasi membuat data baru
		err := tx.Create(&entities).Error
		if err != nil {
			// Check for duplicate entry error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If it's a duplicate entry error, panic duplicate
				return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			// For other errors, return the error
			return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	} else {
		// Jika IncentiveMasterId != 0, ini adalah operasi memperbarui data yang sudah ada
		err := tx.Model(&masterentities.IncentiveMaster{}).
			Where("incentive_level_id = ?", request.IncentiveLevelId).
			Updates(entities).Error
		if err != nil {
			return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities, nil
}

func (r *IncentiveMasterRepositoryImpl) UpdateIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest, Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	entities := masterentities.IncentiveMaster{
		IncentiveLevelId:      request.IncentiveLevelId,
		IncentiveLevelCode:    request.IncentiveLevelCode,
		JobPositionId:         request.JobPositionId,
		IncentiveLevelPercent: request.IncentiveLevelPercent,
		IsActive:              request.IsActive,
	}

	err := tx.Model(&masterentities.IncentiveMaster{}).
		Where("incentive_level_id = ?", Id).
		Updates(entities).Error

	if err != nil {
		return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *IncentiveMasterRepositoryImpl) ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (masterentities.IncentiveMaster, *exceptions.BaseErrorResponse) {
	var entities masterentities.IncentiveMaster

	result := tx.Model(&entities).
		Where("incentive_level_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("incentive with ID %d not found", Id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.IsActive = !entities.IsActive

	// Simpan perubahan
	result = tx.Save(&entities)
	if result.Error != nil {
		return masterentities.IncentiveMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}
