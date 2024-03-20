package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncentiveMasterRepositoryImpl struct {
}

func StartIncentiveMasterRepositoryImpl() masterrepository.IncentiveMasterRepository {
	return &IncentiveMasterRepositoryImpl{}
}

func (r *IncentiveMasterRepositoryImpl) GetAllIncentiveMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masterpayloads.IncentiveMasterListResponse
	var getJobPositionResponse []masterpayloads.JobPositionResponse
	var c *gin.Context
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
		return nil, 0, 0, err
	}

	defer rows.Close()

	if len(responses) == 0 {
		return nil, 0, 0, gorm.ErrRecordNotFound
	}

	jobPositionUrl := "http://10.1.32.26:8000/general-service/api/general/job-position?job_position_id=" + jobPositionId

	errUrlIncentiveMaster := utils.Get(c, jobPositionUrl, &getJobPositionResponse, nil)

	if errUrlIncentiveMaster != nil {
		return nil, 0, 0, errUrlIncentiveMaster
	}

	joinedData := utils.DataFrameInnerJoin(responses, getJobPositionResponse, "JobPositionId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *IncentiveMasterRepositoryImpl) GetIncentiveMasterById(tx *gorm.DB, Id int) (masterpayloads.IncentiveMasterResponse, error) {
	entities := masterentities.IncentiveMaster{}
	response := masterpayloads.IncentiveMasterResponse{}

	err := tx.Model(&entities).
		Where(masterentities.IncentiveMaster{
			IncentiveLevelId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *IncentiveMasterRepositoryImpl) SaveIncentiveMaster(tx *gorm.DB, request masterpayloads.IncentiveMasterRequest) (bool, error) {
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
			return false, err
		}
	} else {
		// Jika IncentiveMasterId != 0, ini adalah operasi memperbarui data yang sudah ada
		err := tx.Model(&masterentities.IncentiveMaster{}).
			Where("incentive_level_id = ?", request.IncentiveLevelId).
			Updates(entities).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *IncentiveMasterRepositoryImpl) ChangeStatusIncentiveMaster(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.IncentiveMaster

	result := tx.Model(&entities).
		Where("incentive_level_id = ?", Id).
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
