package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	// "after-sales/api/exceptions"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"gorm.io/gorm"
)

type IncentiveMasterRepositoryImpl struct {
}

func StartIncentiveMasterRepositoryImpl() masterrepository.IncentiveMasterRepository {
	return &IncentiveMasterRepositoryImpl{}
}

func (r *IncentiveMasterRepositoryImpl) GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masterpayloads.IncentiveMasterListResponse
	var getJobPositionResponse []masterpayloads.JobPositionResponse
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var jobPositionId string
	responseStruct := reflect.TypeOf(masterpayloads.IncentiveMasterListResponse{})

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
		jobPositionId = externalServiceFilter[i].ColumnValue
	}

	// define table struct
	tableStruct := masterpayloads.IncentiveMasterListResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)
	//apply pagination and execute
	rows, err := whereQuery.Scan(&responses).Rows()

	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		// notFoundErr := exceptions.NewNotFoundError("No data found")
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	defer rows.Close()

	jobPositionUrl := config.EnvConfigs.GeneralServiceUrl + "job-position?job_position_id=" + jobPositionId

	errUrlIncentiveMaster := utils.Get(jobPositionUrl, &getJobPositionResponse, nil)

	if errUrlIncentiveMaster != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	joinedData := utils.DataFrameInnerJoin(responses, getJobPositionResponse, "JobPositionId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *IncentiveMasterRepositoryImpl) GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.IncentiveMaster{}
	response := masterpayloads.IncentiveMasterResponse{}

	err := tx.Model(&entities).
		Where(masterentities.IncentiveMaster{
			IncentiveLevelId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *IncentiveMasterRepositoryImpl) SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.IncentiveMaster{
		IncentiveLevelId:      request.IncentiveLevelId,
		IncentiveLevelCode:    request.IncentiveLevelCode,
		JobPositionId:         request.JobPositionId,
		IncentiveLevelPercent: request.IncentiveLevelPercent,
	}

	if request.IncentiveLevelId == 0 {
		// Jika IncentiveMasterId == 0, ini adalah operasi membuat data baru
		err := tx.Create(&entities).Error
		if err != nil {
			// Check for duplicate entry error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If it's a duplicate entry error, panic duplicate
				return false, &exceptionsss_test.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			// For other errors, return the error
			return false, &exceptionsss_test.BaseErrorResponse{
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
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *IncentiveMasterRepositoryImpl) ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (masterentities.IncentiveMaster, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.IncentiveMaster

	result := tx.Model(&entities).
		Where("incentive_level_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masterentities.IncentiveMaster{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("incentive with ID %d not found", Id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return masterentities.IncentiveMaster{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entities.IsActive = !entities.IsActive

	// Simpan perubahan
	result = tx.Save(&entities)
	if result.Error != nil {
		return masterentities.IncentiveMaster{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}
