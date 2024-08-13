package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SupplySlipRepositoryImpl struct {
}

func StartSupplySlipRepositoryImpl() transactionsparepartrepository.SupplySlipRepository {
	return &SupplySlipRepositoryImpl{}
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipById(tx *gorm.DB, Id int, pagination pagination.Pagination) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlip{}
	response := transactionsparepartpayloads.SupplySlipResponse{}
	var getApprovalStatusResponse transactionsparepartpayloads.ApprovalStatusResponse
	var getSupplyTypeResponse transactionsparepartpayloads.SupplyTypeResponse
	var getCustomerResponse transactionsparepartpayloads.CustomerResponse
	var getTechnicianResponse transactionsparepartpayloads.TechnicianResponse
	var getBrandResponse transactionsparepartpayloads.BrandResponse
	var getModelResponse transactionsparepartpayloads.ModelResponse
	var getVariantResponse transactionsparepartpayloads.VariantResponse

	rows, err := tx.Model(&entities).
		Where("supply_system_number = ?", Id).
		Joins("JOIN trx_work_order on trx_supply_slip.work_order_system_number = trx_work_order.work_order_system_number").
		Joins("JOIN mtr_campaign on trx_supply_slip.campaign_id = mtr_campaign.campaign_id").
		Select("trx_supply_slip.*, trx_work_order.*, mtr_campaign.campaign_id, mtr_campaign.campaign_code").
		First(&response).
		Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	// Fetch approval status data
	approvalStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status/" + strconv.Itoa(response.SupplyStatusId)
	errUrlApprovalStatus := utils.Get(approvalStatusUrl, &getApprovalStatusResponse, nil)
	if errUrlApprovalStatus != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlApprovalStatus,
		}
	}
	// Perform inner join with approval status data
	joinedData1 := utils.DataFrameInnerJoin([]transactionsparepartpayloads.SupplySlipResponse{response}, []transactionsparepartpayloads.ApprovalStatusResponse{getApprovalStatusResponse}, "SupplyStatusId")

	// Fetch supply type data
	supplyTypeUrl := config.EnvConfigs.GeneralServiceUrl + "supply-type/" + strconv.Itoa(response.SupplyTypeId)
	errUrlSupplyType := utils.Get(supplyTypeUrl, &getSupplyTypeResponse, nil)
	if errUrlSupplyType != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlSupplyType,
		}
	}
	// Perform inner join with supply type data
	joinedData2 := utils.DataFrameInnerJoin(joinedData1, []transactionsparepartpayloads.SupplyTypeResponse{getSupplyTypeResponse}, "SupplyTypeId")

	// Fetch customer data
	customerUrl := config.EnvConfigs.GeneralServiceUrl + "customer/" + strconv.Itoa(response.CustomerId)
	errUrlCustomer := utils.Get(customerUrl, &getCustomerResponse, nil)
	if errUrlCustomer != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlCustomer,
		}
	}
	// Perform inner join with customer data
	joinedData3 := utils.DataFrameInnerJoin(joinedData2, []transactionsparepartpayloads.CustomerResponse{getCustomerResponse}, "CustomerId")

	// Fetch technician data
	technicianUrl := config.EnvConfigs.GeneralServiceUrl + "user-details-name-and-nickname/" + strconv.Itoa(response.TechnicianId)
	errUrlTechnician := utils.Get(technicianUrl, &getTechnicianResponse, nil)
	if errUrlTechnician != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlTechnician,
		}
	}
	// Perform inner join with technician data
	joinedData4 := utils.DataFrameInnerJoin(joinedData3, []transactionsparepartpayloads.TechnicianResponse{getTechnicianResponse}, "TechnicianId")

	// Fetch brand data
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(response.BrandId)
	errUrlBrand := utils.Get(brandUrl, &getBrandResponse, nil)
	if errUrlBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlBrand,
		}
	}
	// Perform inner join with brand data
	joinedData5 := utils.DataFrameInnerJoin(joinedData4, []transactionsparepartpayloads.BrandResponse{getBrandResponse}, "BrandId")

	// Fetch model data
	modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(response.ModelId)
	errUrlModel := utils.Get(modelUrl, &getModelResponse, nil)
	if errUrlModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlModel,
		}
	}
	// Perform inner join with model data
	joinedData6 := utils.DataFrameInnerJoin(joinedData5, []transactionsparepartpayloads.ModelResponse{getModelResponse}, "ModelId")

	// Fetch variant data
	variantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(response.VariantId)
	errUrlVariant := utils.Get(variantUrl, &getVariantResponse, nil)
	if errUrlVariant != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlVariant,
		}
	}
	// Perform inner join with variant data
	joinedData7 := utils.DataFrameInnerJoin(joinedData6, []transactionsparepartpayloads.VariantResponse{getVariantResponse}, "VariantId")

	result := joinedData7[0]

	// Fetch supply slip details with pagination
	var supplySlipDetail []transactionsparepartpayloads.SupplySlipDetailByHeaderIdResponse
	query := tx.Model(&transactionsparepartentities.SupplySlipDetail{}).
		Where("supply_system_number = ?", Id).
		Joins("JOIN work_order_operation on trx_supply_slip_detail.work_order_operation_id = work_order_operation.work_order_operation_id").
		Joins("JOIN mtr_operation_model_mapping on work_order_operation.operation_id = mtr_operation_model_mapping.operation_model_mapping_id").
		Joins("JOIN mtr_operation_code on mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id").
		Joins("JOIN trx_work_order_item on trx_supply_slip_detail.work_order_item_id = trx_work_order_item.work_order_item_id").
		Joins("JOIN mtr_item on trx_work_order_item.item_id = mtr_item.item_id").
		Joins("JOIN mtr_uom on trx_supply_slip_detail.unit_of_measurement_id = mtr_uom.uom_id").
		Joins("JOIN mtr_warehouse_group on trx_supply_slip_detail.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("JOIN mtr_warehouse_master on trx_supply_slip_detail.warehouse_id = mtr_warehouse_master.warehouse_id").
		Joins("JOIN mtr_warehouse_location on trx_supply_slip_detail.location_id = mtr_warehouse_location.warehouse_location_id").
		Select("trx_supply_slip_detail.*, mtr_operation_code.operation_code, mtr_item.item_code, mtr_item.item_name, mtr_uom.uom_code, mtr_warehouse_group.warehouse_group_code, mtr_warehouse_master.warehouse_code, mtr_warehouse_location.warehouse_location_code").
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errSupplySlipDetail := query.Find(&supplySlipDetail).Error
	if errSupplySlipDetail != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve supply slip details from the database",
			Err:        errSupplySlipDetail,
		}
	}

	responsePayloads := map[string]interface{}{
		"supply_system_number":       result["SupplySystemNumber"],
		"supply_status_id":           result["SupplyStatusId"],
		"supply_status_description":  result["SupplyStatusDescription"],
		"supply_document_number":     result["SupplyDocumentNumber"],
		"supply_date":                result["SupplyDate"],
		"supply_type_id":             result["SupplyTypeId"],
		"supply_type_description":    result["SupplyTypeDescription"],
		"work_order_system_number":   result["WorkOrderSystemNumber"],
		"work_order_document_number": result["WorkOrderDocumentNumber"],
		"work_order_date":            result["WorkOrderDate"],
		"customer_id":                result["CustomerId"],
		"customer_code":              result["CustomerCode"],
		"customer_name":              result["CustomerName"],
		"technician_id":              result["TechnicianId"],
		"technician_name":            result["TechnicianName"],
		"brand_id":                   result["BrandId"],
		"brand_code":                 result["BrandCode"],
		"model_id":                   result["ModelId"],
		"model_code_name":            fmt.Sprintf("%s - %s", result["ModelCode"], result["ModelDescription"]),
		"variant_id":                 result["VariantId"],
		"variant_code_name":          fmt.Sprintf("%s - %s", result["VariantCode"], result["VariantDescription"]),
		"production_year":            result["ProductionYear"],
		"campaign_id":                result["CampaignId"],
		"campaign_code":              result["CampaignCode"],
		"supply_slip_detail": transactionsparepartpayloads.SupplySlipDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows), // Convert int64 to int
			Data:       supplySlipDetail,
		},
	}

	return responsePayloads, nil
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipDetail{}
	response := transactionsparepartpayloads.SupplySlipDetailResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_detail_system_number = ?", Id).
		Joins("JOIN work_order_operation on trx_supply_slip_detail.work_order_operation_id = work_order_operation.work_order_operation_id").
		Joins("JOIN mtr_operation_model_mapping on work_order_operation.operation_id = mtr_operation_model_mapping.operation_model_mapping_id").
		Joins("JOIN mtr_operation_code on mtr_operation_model_mapping.operation_id = mtr_operation_code.operation_id").
		Joins("JOIN trx_work_order_item on trx_supply_slip_detail.work_order_item_id = trx_work_order_item.work_order_item_id").
		Joins("JOIN mtr_item on trx_work_order_item.item_id = mtr_item.item_id").
		Joins("JOIN mtr_uom on trx_supply_slip_detail.unit_of_measurement_id = mtr_uom.uom_id").
		Joins("JOIN mtr_warehouse_group on trx_supply_slip_detail.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("JOIN mtr_warehouse_master on trx_supply_slip_detail.warehouse_id = mtr_warehouse_master.warehouse_id").
		Joins("JOIN mtr_warehouse_location on trx_supply_slip_detail.location_id = mtr_warehouse_location.warehouse_location_id").
		Select("trx_supply_slip_detail.*, mtr_operation_code.operation_code, mtr_item.item_code, mtr_item.item_name, mtr_uom.uom_code, mtr_warehouse_group.warehouse_group_code, mtr_warehouse_master.warehouse_code, mtr_warehouse_location.warehouse_location_code").
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

func (r *SupplySlipRepositoryImpl) SaveSupplySlip(tx *gorm.DB, request transactionsparepartentities.SupplySlip) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlip{
		SupplyStatusId:        request.SupplyStatusId,
		SupplyDate:            request.SupplyDate,
		SupplyDocumentNumber:  " ",
		SupplyTypeId:          request.SupplyTypeId,
		CompanyId:             request.CompanyId,
		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
		TechnicianId:          request.TechnicianId,
		CampaignId:            request.CampaignId,
		Remark:                request.Remark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlip{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipRepositoryImpl) SaveSupplySlipDetail(tx *gorm.DB, request transactionsparepartentities.SupplySlipDetail) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse) {
	total := request.QuantitySupply - 0
	entities := transactionsparepartentities.SupplySlipDetail{
		SupplySystemNumber:   request.SupplySystemNumber,
		WorkOrderOperationId: request.WorkOrderOperationId,
		WorkOrderItemId:      request.WorkOrderItemId,
		LocationId:           request.LocationId,
		UnitOfMeasurementId:  request.UnitOfMeasurementId,
		QuantitySupply:       request.QuantitySupply,
		QuantityReturn:       0,
		QuantityDemand:       request.QuantityDemand,
		CostOfGoodsSold:      0,
		WorkOrderDetailId:    request.WorkOrderDetailId,
		WarehouseGroupId:     request.WarehouseGroupId,
		WarehouseId:          request.WarehouseId,
		QuantityTotal:        total,
		PR:                   request.PR,
		QuantityPR:           request.QuantityPR,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlipDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipRepositoryImpl) GetAllSupplySlip(tx *gorm.DB, internalFilter []utils.FilterCondition, externalFilter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.SupplySlipSearchResponse
	var getSupplyTypeResponse transactionsparepartpayloads.SupplyTypeResponse
	var getSupplyTypeAllResponse []transactionsparepartpayloads.SupplyTypeResponse
	var getApprovalStatusResponse transactionsparepartpayloads.ApprovalStatusResponse
	var getApprovalStatusAllResponse []transactionsparepartpayloads.ApprovalStatusResponse

	supplyTypeId := ""
	approvalStatusId := ""

	// apply external services filter

	for i := 0; i < len(externalFilter); i++ {
		if strings.Contains(externalFilter[i].ColumnField, "supply_type_id") {
			supplyTypeId = externalFilter[i].ColumnValue
		} else if strings.Contains(externalFilter[i].ColumnField, "approval_status_id") {
			approvalStatusId = externalFilter[i].ColumnValue
		}
	}

	// define table struct
	tableStruct := transactionsparepartpayloads.SupplySlipSearchResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatementTransaction(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilterSearch(joinTable, internalFilter)

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
		var response transactionsparepartpayloads.SupplySlipSearchResponse
		if err := rows.Scan(&response.SupplySystemNumber, &response.SupplyDocumentNumber, &response.SupplyDate, &response.SupplyTypeId, &response.WorkOrderSystemNumber, &response.WorkOrderDocumentNumber, &response.CustomerId, &response.SupplyStatusId); err != nil {
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

	// Fetch supply type data
	if supplyTypeId != "" {
		supplyTypeUrl := config.EnvConfigs.GeneralServiceUrl + "supply-type/" + supplyTypeId
		errUrlSupplyType := utils.Get(supplyTypeUrl, &getSupplyTypeResponse, nil)
		if errUrlSupplyType != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlSupplyType,
			}
		}
		// Perform inner join with supply type data
		joinedData1 = utils.DataFrameInnerJoin(responses, []transactionsparepartpayloads.SupplyTypeResponse{getSupplyTypeResponse}, "SupplyTypeId")
	} else {
		supplyTypeUrl := config.EnvConfigs.GeneralServiceUrl + "supply-type"
		errUrlSupplyType := utils.Get(supplyTypeUrl, &getSupplyTypeAllResponse, nil)
		if errUrlSupplyType != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUrlSupplyType,
			}
		}
		// Perform inner join with supply type data
		joinedData1 = utils.DataFrameInnerJoin(responses, getSupplyTypeAllResponse, "SupplyTypeId")
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
		joinedData2 = utils.DataFrameInnerJoin(joinedData1, []transactionsparepartpayloads.ApprovalStatusResponse{getApprovalStatusResponse}, "SupplyStatusId")
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
		joinedData2 = utils.DataFrameInnerJoin(joinedData1, getApprovalStatusAllResponse, "SupplyStatusId")
	}

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData2, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *SupplySlipRepositoryImpl) UpdateSupplySlip(tx *gorm.DB, req transactionsparepartentities.SupplySlip, id int) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse) {
	var entity transactionsparepartentities.SupplySlip

	err := tx.Model(entity).Where(transactionsparepartentities.SupplySlip{SupplySystemNumber: id}).Updates(req).First(&entity).Error
	if err != nil {
		return transactionsparepartentities.SupplySlip{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity, nil
}

func (r *SupplySlipRepositoryImpl) UpdateSupplySlipDetail(tx *gorm.DB, req transactionsparepartentities.SupplySlipDetail, id int) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse) {
	var entity transactionsparepartentities.SupplySlipDetail

	err := tx.Model(entity).Where(transactionsparepartentities.SupplySlipDetail{SupplyDetailSystemNumber: id}).Updates(req).First(&entity).Error
	if err != nil {
		return transactionsparepartentities.SupplySlipDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity, nil
}

func (r *SupplySlipRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, supplySlipId int) (string, *exceptions.BaseErrorResponse) {
	var supplySlip transactionsparepartentities.SupplySlip

	err1 := tx.Model(&transactionsparepartentities.SupplySlip{}).
		Where("supply_system_number = ?", supplySlipId).
		First(&supplySlip).
		Error
	if err1 != nil {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve supply slip from the database: %v", err1)}
	}

	var workOrder transactionworkshopentities.WorkOrder
	var brandResponse transactionworkshoppayloads.BrandDocResponse

	workOrderId := supplySlip.WorkOrderSystemNumber

	// Get the work order based on the work order system number
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("work_order_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrder.BrandId)
	errUrl := utils.Get(brandUrl, &brandResponse, nil)
	if errUrl != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrl,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("SPSS/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.WorkOrderSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.WorkOrderDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.WorkOrderDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPSS/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

func (r *SupplySlipRepositoryImpl) SubmitSupplySlip(tx *gorm.DB, supplySlipId int) (bool, string, *exceptions.BaseErrorResponse) {
	var entity transactionsparepartentities.SupplySlip
	err := tx.Model(&transactionsparepartentities.SupplySlip{}).Where("supply_system_number = ?", supplySlipId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", &exceptions.BaseErrorResponse{Message: "No supply slip data found"}
		}
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve supply slip from the database: %v", err)}
	}

	if entity.SupplyDocumentNumber == " " && entity.SupplyStatusId == 4 {
		//Generate new document number
		newDocumentNumber, genErr := r.GenerateDocumentNumber(tx, entity.SupplySystemNumber)
		if genErr != nil {
			return false, "", genErr
		}
		//newDocumentNumber 

		entity.SupplyDocumentNumber = newDocumentNumber

		// Update work order status to 8 (Wait Approve)
		entity.SupplyStatusId = 8

		err = tx.Save(&entity).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to submit the supply slip: %v", err)}
		}

		return true, newDocumentNumber, nil
	} else {

		return false, "", &exceptions.BaseErrorResponse{Message: "Document number has already been generated"}
	}
}