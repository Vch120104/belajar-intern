package masterrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type AgreementRepositoryImpl struct {
}

func StartAgreementRepositoryImpl() masterrepository.AgreementRepository {
	return &AgreementRepositoryImpl{}
}

func (r *AgreementRepositoryImpl) GetAgreementById(tx *gorm.DB, AgreementId int) (masterpayloads.AgreementResponse, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Agreement{}
	response := masterpayloads.AgreementResponse{}

	err := tx.Model(&entities).
		Where(masterentities.Agreement{
			AgreementId: AgreementId,
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

func (r *AgreementRepositoryImpl) SaveAgreement(tx *gorm.DB, req masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masterentities.Agreement{
		IsActive:    req.IsActive,
		AgreementId: req.AgreementId,
		CustomerId:  req.CustomerId,
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

func (r *AgreementRepositoryImpl) ChangeStatusAgreement(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterentities.Agreement

	result := tx.Model(&entities).
		Where("forecast_master_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

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

func (r *AgreementRepositoryImpl) GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	// Define variables
	var (
		responses             []masterpayloads.AgreementListResponse
		getSupplierResponse   []masterpayloads.SupplierResponse
		getOrderTypeResponse  []masterpayloads.OrderTypeResponse
		internalServiceFilter []utils.FilterCondition
		supplierName          string
		orderTypeName         string
		responseStruct        = reflect.TypeOf(masterpayloads.AgreementListResponse{})
	)

	// Apply internal and external service filters
	for _, fc := range filterCondition {
		var flag bool
		for j := 0; j < responseStruct.NumField(); j++ {
			if fc.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, fc)
				flag = true
				break
			}
		}
		if !flag {
			if strings.Contains(fc.ColumnField, "supplier_name") {
				supplierName = fc.ColumnValue
			} else {
				orderTypeName = fc.ColumnValue
			}
		}
	}

	// Define table struct
	tableStruct := masterpayloads.AgreementListResponse{}

	// Create join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply internal service filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Execute query
	if err := whereQuery.Scan(&responses).Error; err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Check if no records found
	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        gorm.ErrRecordNotFound,
		}
	}

	// Handle supplier and order type filters
	if supplierName != "" || orderTypeName != "" {
		supplierURL := config.EnvConfigs.GeneralServiceUrl + "/api/general/filter-supplier-master?supplier_name=" + supplierName
		if err := utils.Get(supplierURL, &getSupplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		joinedData := utils.DataFrameInnerJoin(responses, getSupplierResponse, "SupplierId")

		orderTypeURL := config.EnvConfigs.GeneralServiceUrl + "/api/general/order-type-filter?order_type_name=" + orderTypeName
		if err := utils.Get(orderTypeURL, &getOrderTypeResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		joinedData = utils.DataFrameInnerJoin(joinedData, getOrderTypeResponse, "OrderTypeId")

		// Paginate data
		dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)
		return dataPaginate, totalPages, totalRows, nil
	}

	// Paginate data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)
	return dataPaginate, totalPages, totalRows, nil
}
