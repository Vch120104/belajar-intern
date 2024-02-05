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

type ForecastMasterRepositoryImpl struct {
}

func StartForecastMasterRepositoryImpl() masterrepository.ForecastMasterRepository {
	return &ForecastMasterRepositoryImpl{}
}

func (r *ForecastMasterRepositoryImpl) GetForecastMasterById(tx *gorm.DB, forecastMasterId int) (masterpayloads.ForecastMasterResponse, error) {
	entities := masterentities.ForecastMaster{}
	response := masterpayloads.ForecastMasterResponse{}

	err := tx.Model(&entities).
		Where(masterentities.ForecastMaster{
			ForecastMasterId: forecastMasterId,
		}).
		First(&response).
		Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *ForecastMasterRepositoryImpl) SaveForecastMaster(tx *gorm.DB, req masterpayloads.ForecastMasterResponse) (bool, error) {
	entities := masterentities.ForecastMaster{
		IsActive:                   req.IsActive,
		ForecastMasterId:           req.ForecastMasterId,
		SupplierId:                 req.SupplierId,
		MovingCodeId:               req.MovingCodeId,
		OrderTypeId:                req.OrderTypeId,
		ForecastMasterLeadTime:     req.ForecastMasterLeadTime,
		ForecastMasterSafetyFactor: req.ForecastMasterSafetyFactor,
		ForecastMasterOrderCycle:   req.ForecastMasterOrderCycle,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ForecastMasterRepositoryImpl) ChangeStatusForecastMaster(tx *gorm.DB, Id int) (bool, error) {
	var entities masterentities.ForecastMaster

	result := tx.Model(&entities).
		Where("forecast_master_id = ?", Id).
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

func (r *ForecastMasterRepositoryImpl) GetAllForecastMaster(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error) {
	var responses []masterpayloads.ForecastMasterListResponse
	var getSupplierResponse []masterpayloads.SupplierResponse
	var getOrderTypeResponse []masterpayloads.OrderTypeResponse
	var c *gin.Context
	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var supplierName string
	var orderTypeName string
	responseStruct := reflect.TypeOf(masterpayloads.ForecastMasterListResponse{})

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
		orderTypeName = externalServiceFilter[i].ColumnValue
	}

	// define table struct
	tableStruct := masterpayloads.ForecastMasterListResponse{}
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

	supplierUrl := "http://10.1.32.26:8000/general-service/api/general/filter-supplier-master?supplier_name=" + supplierName

	errUrlSupplier := utils.Get(c, supplierUrl, &getSupplierResponse, nil)

	if errUrlSupplier != nil {
		return nil, 0, 0, errUrlSupplier
	}

	joinedData := utils.DataFrameInnerJoin(responses, getSupplierResponse, "SupplierId")

	orderTypeUrl := "http://10.1.32.26:8000/general-service/api/general/order-type-filter?order_type_name=" + orderTypeName

	errUrlOrderType := utils.Get(c, orderTypeUrl, &getOrderTypeResponse, nil)

	if errUrlOrderType != nil {
		return nil, 0, 0, errUrlOrderType
	}

	joinedData2 := utils.DataFrameInnerJoin(joinedData, getOrderTypeResponse, "OrderTypeId")

	print(joinedData2)

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)

	return dataPaginate, totalPages, totalRows, nil
}
