package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type WorkOrderBypassRepositoryImpl struct {
}

func OpenWorkOrderBypassRepositoryImpl() transactionworkshoprepository.WorkOrderBypassRepository {
	return &WorkOrderBypassRepositoryImpl{}
}

// uspg_wtWorkOrder2_Select
// IF @Option = 13
//
//	--USE FOR : * SELECT DATA BY EMPLOYEE FOR BYPASS TO QC
//	--USE IN MODUL :
func (r *WorkOrderBypassRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []transactionworkshoppayloads.WorkOrderDetailBypassRequest

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.WorkOrderDetailBypassRequest{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	whereQuery = whereQuery.Where("work_order_system_number > 0 and line_type_id = 1 and service_status_id IN (?,?,?,?,?)", utils.SrvStatDraft, utils.SrvStatStart, utils.SrvStatPending, utils.SrvStatStop, utils.SrvStatReOrder)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, payload := range payloads {
		// Fetch work order data from internal services
		ModelURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(payload.WorkOrderSystemNumber)
		var getModelResponse transactionworkshoppayloads.WorkOrderResponse
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch work order data from external service",
				Err:        err,
			}
		}

		// Fetch item data from internal services
		ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(payload.ItemId)
		var getItemResponse transactionworkshoppayloads.ItemServiceRequestDetail
		if err := utils.Get(ItemURL, &getItemResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch item data from external service",
				Err:        err,
			}
		}

		// Fetch line type from internal services
		getOperationResponse, lineErr := generalserviceapiutils.GetLineTypeById(payload.LineTypeId)
		if lineErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch line type data from external service",
				Err:        lineErr.Err,
			}
		}

		// Fetch service status from internal services
		ServiceStatusURL := config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(payload.ServiceStatusId)
		var getServiceStatusResponse transactionworkshoppayloads.ServiceStatusResponse
		if err := utils.Get(ServiceStatusURL, &getServiceStatusResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch service status data from external service",
				Err:        err,
			}
		}

		result := map[string]interface{}{
			"work_order_system_number":   payload.WorkOrderSystemNumber,
			"work_order_document_number": getModelResponse.WorkOrderDocumentNumber,
			"line_type_id":               payload.LineTypeId,
			"line_type_name":             getOperationResponse.LineTypeName,
			"item_id":                    payload.ItemId,
			"item_code":                  getItemResponse.ItemCode,
			"item_name":                  getItemResponse.ItemName,
			"frt_quantity":               payload.FrtQuantity,
			"service_status_name":        getServiceStatusResponse.ServiceStatusName,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *WorkOrderBypassRepositoryImpl) GetById(tx *gorm.DB, id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse) {
	var workOrderResponse transactionworkshoppayloads.WorkOrderBypassResponse
	var tableStruct transactionworkshoppayloads.WorkOrderDetailBypassRequest

	// Create join query
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := joinTable.Where("work_order_system_number = ?", id)

	// Execute the query and populate tableStruct
	if err := whereQuery.First(&tableStruct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderBypassResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Record not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.WorkOrderBypassResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch data work order from internal services
	ModelURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(tableStruct.WorkOrderSystemNumber)
	//fmt.Println("Fetching  work order data from:", ModelURL)
	var getModelResponse transactionworkshoppayloads.WorkOrderResponse
	if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch model data from external service",
			Err:        err,
		}
	}

	// fetch data item from internal services
	ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(tableStruct.ItemId)
	//fmt.Println("Fetching  item data from:", ItemURL)
	var getItemResponse transactionworkshoppayloads.ItemServiceRequestDetail
	if err := utils.Get(ItemURL, &getItemResponse, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch item data from external service",
			Err:        err,
		}
	}

	// fetch data operation from internal services
	getOperationResponse, lineErr := generalserviceapiutils.GetLineTypeById(tableStruct.LineTypeId)
	if lineErr != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch line type data from external service",
			Err:        lineErr.Err,
		}
	}

	// Map the data to the response struct
	workOrderResponse = transactionworkshoppayloads.WorkOrderBypassResponse{
		WorkOrderSystemNumber:   tableStruct.WorkOrderSystemNumber,
		WorkOrderDocumentNumber: getModelResponse.WorkOrderDocumentNumber,
		LineTypeId:              tableStruct.LineTypeId,
		LineTypeName:            getOperationResponse.LineTypeName,
		ItemId:                  tableStruct.ItemId,
		ItemName:                getItemResponse.ItemName,
	}

	return workOrderResponse, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 8
// --USE IN MODUL : * UPDATE DATA BY KEY -- BYPASS OPERATION TO QC PASS
func (r *WorkOrderBypassRepositoryImpl) Bypass(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshoppayloads.WorkOrderBypassResponseDetail, *exceptions.BaseErrorResponse) {
	var wo transactionworkshopentities.WorkOrder
	var carWash transactionjpcbentities.CarWash
	var count int64
	var lineTypeOperation = 1

	// Retrieve WorkOrder and WorkOrderDetail
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", id).
		First(&wo).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Work Order not found",
			Err:        err,
		}
	}

	if wo.WorkOrderDocumentNumber == "" {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Work Order Document Number must be filled before Bypass",
			Err:        errors.New("work order document number must be filled before bypass"),
		}
	}

	// Update WorkOrderDetail
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_operation_item_line = ?", id, lineTypeOperation).
		Updates(map[string]interface{}{
			"service_status_id":             utils.SrvStatQcPass,
			"quality_control_pass_datetime": time.Now(),
			"bypass":                        true,
			"technician_id":                 request.TechnicianId,
		}).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update Work Order Detail",
			Err:        err,
		}
	}

	// Delete from WorkOrderTechAlloc
	if err := tx.Where("work_order_system_number = ? AND work_order_line = ?", id, lineTypeOperation).
		Delete(&transactionworkshopentities.WorkOrderAllocation{}).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete Work Order Allocation",
			Err:        err,
		}
	}

	// Delete from ServiceLog
	if err := tx.Where("work_order_system_number = ? AND work_order_line = ?", id, lineTypeOperation).
		Delete(&transactionworkshopentities.ServiceLog{}).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete Service Log",
			Err:        err,
		}
	}

	// Update WorkOrder status based on conditions
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND line_type_id = ? AND service_status_id <> ?", id, lineTypeOperation, utils.SrvStatQcPass).
		Count(&count).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count Work Order Detail",
			Err:        err,
		}
	}

	if count == 0 {
		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Updates(map[string]interface{}{
				"work_order_status_id": utils.WoStatQC,
			}).Error; err != nil {
			return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update Work Order status",
				Err:        err,
			}
		}
	}

	// Insert into CarWash if not exists
	if err := tx.Model(&transactionjpcbentities.CarWash{}).
		Where("work_order_system_number = ?", id).
		First(&carWash).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Car Wash",
			Err:        err,
		}
	}

	if carWash.WorkOrderSystemNumber == 0 {
		carWash = transactionjpcbentities.CarWash{
			CompanyId:             1,
			WorkOrderSystemNumber: id,
			CarWashDate:           time.Now(),
			StartTime:             0,
			EndTime:               0,
			ActualTime:            0,
		}

		if err := tx.Create(&carWash).Error; err != nil {
			return transactionworkshoppayloads.WorkOrderBypassResponseDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create Car Wash",
				Err:        err,
			}
		}
	}

	// Prepare the response
	response := transactionworkshoppayloads.WorkOrderBypassResponseDetail{
		WorkOrderSystemNumber: wo.WorkOrderSystemNumber,
		ServiceStatusId:       utils.SrvStatQcPass,
		TechnicianId:          request.TechnicianId,
	}

	return response, nil
}
