package masterwarehouserepositoryimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	"after-sales/api/config"
	"after-sales/api/exceptions"
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	// utils "after-sales/api/utils"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationDefinitionRepositoryImpl struct {
}

func OpenWarehouseLocationDefinitionImpl() masterwarehouserepository.WarehouseLocationDefinitionRepository {
	return &WarehouseLocationDefinitionRepositoryImpl{}
}

func (r *WarehouseLocationDefinitionRepositoryImpl) Save(tx *gorm.DB, request masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptionsss_test.BaseErrorResponse) {

	var warehouseMaster = masterwarehouseentities.WarehouseLocationDefinition{
		IsActive:                               request.IsActive,
		WarehouseLocationDefinitionId:          request.WarehouseLocationDefinitionId,
		WarehouseLocationDefinitionLevelId:     request.WarehouseLocationDefinitionLevelId,
		WarehouseLocationDefinitionLevelCode:   request.WarehouseLocationDefinitionLevelCode,
		WarehouseLocationDefinitionDescription: request.WarehouseLocationDefinitionDescription,
	}

	err := tx.Save(&warehouseMaster).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *WarehouseLocationDefinitionRepositoryImpl) GetById(tx *gorm.DB, Id int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptionsss_test.BaseErrorResponse) {

	entities := masterwarehouseentities.WarehouseLocationDefinition{}
	response := masterwarehousepayloads.WarehouseLocationDefinitionResponse{}

	err := tx.Model(&entities).
		Where("warehouse_location_definition_id = ?", Id).
		First(&response).
		Error

	if err != nil {
		return masterwarehousepayloads.WarehouseLocationDefinitionResponse{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return response, nil
}

func (r *WarehouseLocationDefinitionRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masterwarehousepayloads.WarehouseLocationDefinitionResponse
	var getWhLevelResponse []masterwarehousepayloads.WarehouseLocationDefinitionLevelResponse
	var internalServiceFilter []utils.FilterCondition
	var warehouseLocationDefinitionLevelId int

	responseStruct := reflect.TypeOf(masterwarehousepayloads.WarehouseLocationDefinitionResponse{})

	// Filter internal service conditions
	for _, condition := range filterCondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	// Apply internal service filter conditions
	tableStruct := masterwarehousepayloads.WarehouseLocationDefinitionRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Fetch data from database
	err := whereQuery.Scan(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		notFoundErr := exceptions.NewNotFoundError("No data found")
		panic(notFoundErr)
	}

	// Extract warehouse location definition level ID from the first response
	warehouseLocationDefinitionLevelId = responses[0].WarehouseLocationDefinitionLevelId

	// Fetch warehouse location definition level data from external service
	whLevelUrl := config.EnvConfigs.AfterSalesServiceUrl + "/warehouse-location-definition/popup-level?warehouse_location_definition_level_id=" + strconv.Itoa(warehouseLocationDefinitionLevelId)
	err = utils.Get(whLevelUrl, &getWhLevelResponse, nil)
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *WarehouseLocationDefinitionRepositoryImpl) ChangeStatus(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entity masterwarehouseentities.WarehouseLocationDefinition

	// Cari entitas berdasarkan ID
	result := tx.Model(&entity).
		Where("warehouse_location_definition_id = ?", Id).
		First(&entity)

	// Periksa apakah entitas ditemukan
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("warehouse Loc with ID %d not found", Id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	// Ubah status entitas
	entity.IsActive = !entity.IsActive

	// Simpan perubahan
	result = tx.Save(&entity)
	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *WarehouseLocationDefinitionRepositoryImpl) PopupWarehouseLocationLevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masterwarehousepayloads.WarehouseLocationDefinitionLevelResponse

	// Fetch data from database with joins and conditions
	query := tx.Table("mtr_warehouse_location_definition_level")

	// Apply filter conditions
	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	err := query.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		notFoundErr := exceptions.NewNotFoundError("No data found")
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        notFoundErr,
		}
	}

	// Perform pagination
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}
