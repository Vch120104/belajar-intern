package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type SupplySlipReturnRepositoryImpl struct {
}

func StartSupplySlipReturnRepositoryImpl() transactionsparepartrepository.SupplySlipReturnRepository {
	return &SupplySlipReturnRepositoryImpl{}
}

func (r *SupplySlipReturnRepositoryImpl) SaveSupplySlipReturn(tx *gorm.DB, request transactionsparepartentities.SupplySlipReturn) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturn{
		SupplyID:                   request.SupplyID,
		SupplyReturnDate:           request.SupplyReturnDate,
		SupplyReturnDocumentNumber: " ",
		SupplyReturnStatusId:       request.SupplyReturnStatusId,
		Remark:                     request.Remark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlipReturn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipReturnRepositoryImpl) SaveSupplySlipReturnDetail(tx *gorm.DB, request transactionsparepartentities.SupplySlipReturnDetail) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturnDetail{
		SupplyReturnID:        request.SupplyReturnID,
		SupplyDetailID:        request.SupplyDetailID,
		QuantityReturn:        request.QuantityReturn,
		SupplyReturnReasonID:  request.SupplyReturnReasonID,
		CostOfGoodsSoldReturn: 0,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlipReturnDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipReturnRepositoryImpl) GetAllSupplySlipReturn(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturn{}
	var responses []transactionsparepartpayloads.SupplySlipReturnSearchResponse
	// var getCustomerResponse transactionsparepartpayloads.CustomerResponse
	var getCustomerAllResponse []transactionsparepartpayloads.CustomerResponse
	var getApprovalStatusResponse transactionsparepartpayloads.SupplyReturnStatusResponse
	var getApprovalStatusAllResponse []transactionsparepartpayloads.SupplyReturnStatusResponse

	customerName := ""
	approvalStatusId := ""

	// apply external services filter

	for i := 0; i < len(externalFilter); i++ {
		if strings.Contains(externalFilter[i].ColumnField, "customer_name") {
			customerName = externalFilter[i].ColumnValue
		} else if strings.Contains(externalFilter[i].ColumnField, "approval_status_id") {
			approvalStatusId = externalFilter[i].ColumnValue
		}
	}

	joinTable := tx.Model(&entities).
		Joins("JOIN trx_supply_slip on trx_supply_slip_return.supply_system_number = trx_supply_slip.supply_system_number").
		Joins("JOIN trx_work_order on trx_supply_slip.work_order_system_number = trx_work_order.work_order_system_number").
		Select("trx_supply_slip_return.supply_return_system_number, trx_supply_slip_return.supply_return_document_number, trx_supply_slip_return.supply_return_date, trx_supply_slip.supply_document_number, trx_work_order.work_order_document_number, trx_work_order.customer_id, trx_supply_slip_return.supply_return_status_id")

	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalFilter)

	// Execute the query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	// Scan the results into the responses slice
	for rows.Next() {
		var response transactionsparepartpayloads.SupplySlipReturnSearchResponse
		if err := rows.Scan(&response.SupplyReturnSystemNumber, &response.SupplyReturnDocumentNumber, &response.SupplyReturnDate, &response.SupplyDocumentNumber, &response.WorkOrderDocumentNumber, &response.CustomerId, &response.SupplyReturnStatusId); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	var joinedData1 []map[string]interface{}

	// Fetch customer data
	if customerName != "" {
		customerUrl := config.EnvConfigs.GeneralServiceUrl + "customer-by-name/" + customerName
		errUrlCustomer := utils.Get(customerUrl, &getCustomerAllResponse, nil)
		if errUrlCustomer != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlCustomer,
			}
		}
		// Perform inner join with customer data
		joinedData1, err = utils.DataFrameInnerJoin(responses, getCustomerAllResponse, "CustomerId")

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	} else {
		customerUrl := config.EnvConfigs.GeneralServiceUrl + "customers?page=0&limit=1000"
		errUrlCustomer := utils.Get(customerUrl, &getCustomerAllResponse, nil)
		if errUrlCustomer != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlCustomer,
			}
		}
		// Perform inner join with customer data
		joinedData1, err = utils.DataFrameInnerJoin(responses, getCustomerAllResponse, "CustomerId")

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	var joinedData2 []map[string]interface{}

	// Fetch approval status data
	if approvalStatusId != "" {
		approvalStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status/" + approvalStatusId
		errUrlapprovalStatus := utils.Get(approvalStatusUrl, &getApprovalStatusResponse, nil)
		if errUrlapprovalStatus != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlapprovalStatus,
			}
		}
		// Perform inner join with supply type data
		joinedData2, err = utils.DataFrameInnerJoin(joinedData1, []transactionsparepartpayloads.SupplyReturnStatusResponse{getApprovalStatusResponse}, "SupplyReturnStatusId")

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	} else {
		approvalStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status"
		errUrlapprovalStatus := utils.Get(approvalStatusUrl, &getApprovalStatusAllResponse, nil)
		if errUrlapprovalStatus != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlapprovalStatus,
			}
		}
		// Perform inner join with supply type data
		joinedData2, err = utils.DataFrameInnerJoin(joinedData1, getApprovalStatusAllResponse, "SupplyReturnStatusId")

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *SupplySlipReturnRepositoryImpl) GetSupplySlipReturnById(tx *gorm.DB, Id int, pagination pagination.Pagination, supplySlip map[string]interface{}) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturn{}
	response := transactionsparepartpayloads.SupplySlipReturnResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_return_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	responsePayloads := map[string]interface{}{
		"supply_return_system_number":   response.SupplyReturnSystemNumber,
		"supply_return_status_id":       response.SupplyReturnStatusId,
		"supply_return_document_number": response.SupplyReturnDocumentNumber,
		"supply_return_date":            response.SupplyReturnDate,
		"supply_system_number":          response.SupplyReturnSystemNumber,
		"supply_document_number":        supplySlip["supply_document_number"],
		"supply_date":                   supplySlip["supply_date"],
		"supply_type_id":                supplySlip["supply_type_id"],
		"supply_type_description":       supplySlip["supply_type_description"],
		"work_order_system_number":      supplySlip["work_order_system_number"],
		"work_order_document_number":    supplySlip["work_order_document_number"],
		"work_order_date":               supplySlip["work_order_date"],
		"customer_id":                   supplySlip["customer_id"],
		"customer_code":                 supplySlip["customer_code"],
		"customer_name":                 supplySlip["customer_name"],
		"technician_id":                 supplySlip["technician_id"],
		"technician_name":               supplySlip["technician_name"],
		"brand_id":                      supplySlip["brand_id"],
		"brand_code":                    supplySlip["brand_code"],
		"model_id":                      supplySlip["model_id"],
		"model_code_name":               supplySlip["model_code_name"],
		"variant_id":                    supplySlip["variant_id"],
		"variant_code_name":             supplySlip["variant_code_name"],
		"production_year":               supplySlip["production_year"],
		"campaign_id":                   supplySlip["campaign_id"],
		"campaign_code":                 supplySlip["campaign_code"],
	}

	return responsePayloads, nil

}

func (r *SupplySlipReturnRepositoryImpl) GetSupplySlipReturnDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipReturnDetailResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturnDetail{}
	response := transactionsparepartpayloads.SupplySlipReturnDetailResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_return_detail_system_number = ?", Id).
		Joins("JOIN trx_supply_slip_detail on trx_supply_slip_return_detail.supply_detail_system_number = trx_supply_slip_detail.supply_detail_system_number").
		Joins("JOIN work_order_operation on trx_supply_slip_detail.work_order_operation_id = work_order_operation.work_order_operation_id").
		Joins("JOIN mtr_operation_model_mapping on work_order_operation.operation_id = mtr_operation_model_mapping.operation_model_mapping_id").
		Joins("JOIN mtr_operation_code on mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id").
		Joins("JOIN trx_work_order_item on trx_supply_slip_detail.work_order_item_id = trx_work_order_item.work_order_item_id").
		Joins("JOIN mtr_item on trx_work_order_item.item_id = mtr_item.item_id").
		Joins("JOIN mtr_uom on trx_supply_slip_detail.unit_of_measurement_id = mtr_uom.uom_id").
		Joins("JOIN mtr_warehouse_group on trx_supply_slip_detail.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("JOIN mtr_warehouse_master on trx_supply_slip_detail.warehouse_id = mtr_warehouse_master.warehouse_id").
		Joins("JOIN mtr_warehouse_location on trx_supply_slip_detail.location_id = mtr_warehouse_location.warehouse_location_id").
		Select("trx_supply_slip_return_detail.*, mtr_operation_code.operation_code, mtr_item.item_code, mtr_item.item_name, mtr_uom.uom_code, mtr_warehouse_group.warehouse_group_code, mtr_warehouse_master.warehouse_code, mtr_warehouse_location.warehouse_location_code").
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *SupplySlipReturnRepositoryImpl) UpdateSupplySlipReturn(tx *gorm.DB, req transactionsparepartentities.SupplySlipReturn, id int) (transactionsparepartentities.SupplySlipReturn, *exceptions.BaseErrorResponse) {
	var entity transactionsparepartentities.SupplySlipReturn

	err := tx.Model(entity).Where(transactionsparepartentities.SupplySlipReturn{SupplyReturnSystemNumber: id}).Updates(req).First(&entity).Error
	if err != nil {
		return transactionsparepartentities.SupplySlipReturn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity, nil
}

func (r *SupplySlipReturnRepositoryImpl) UpdateSupplySlipReturnDetail(tx *gorm.DB, req transactionsparepartentities.SupplySlipReturnDetail, id int) (transactionsparepartentities.SupplySlipReturnDetail, *exceptions.BaseErrorResponse) {
	var entity transactionsparepartentities.SupplySlipReturnDetail

	err := tx.Model(entity).Where(transactionsparepartentities.SupplySlipReturnDetail{SupplyReturnDetailSystemNumber: id}).Updates(req).First(&entity).Error
	if err != nil {
		return transactionsparepartentities.SupplySlipReturnDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity, nil
}

func (r *SupplySlipReturnRepositoryImpl) GetSupplySlipId(tx *gorm.DB, Id int) (int, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipReturn{}
	response := transactionsparepartpayloads.SupplySlipReturnResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_return_system_number = ?", Id).
		Select("trx_supply_slip_return.*").
		First(&response).
		Rows()

	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response.SupplySystemNumber, nil

}
