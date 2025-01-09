package transactionbodyshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionbodyshoprepository "after-sales/api/repositories/transaction/bodyshop"
	"after-sales/api/utils"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type QualityControlBodyshopRepositoryImpl struct {
}

func OpenQualityControlBodyshopRepositoryImpl() transactionbodyshoprepository.QualityControlBodyshopRepository {
	return &QualityControlBodyshopRepositoryImpl{}
}

// uspg_wtWorkOrder0_Select
// IF @Option = 7
// USE IN MODUL : AWS - 006 QUALITY CONTROL PAGE 1 REQ: ???
func (r *QualityControlBodyshopRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var entities []transactionbodyshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, transactionbodyshoppayloads.QualityControlRequest{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_status_id = ?", utils.WoStatStop) // 40 Stop

	if err := whereQuery.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var convertedResponses []transactionbodyshoppayloads.QualityControlResponse

	for _, entity := range entities {
		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
		if modelErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve model data from the external API",
				Err:        modelErr.Err,
			}
		}

		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
		if variantErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve variant data from the external API",
				Err:        variantErr.Err,
			}
		}

		vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
		if vehicleErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        vehicleErr.Err,
			}
		}

		CustomerUrl := config.EnvConfigs.SalesServiceUrl + "customer/" + strconv.Itoa(entity.CustomerId)
		var customerResponses transactionworkshoppayloads.CustomerResponse
		errCustomer := utils.Get(CustomerUrl, &customerResponses, nil)
		if errCustomer != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve customer data from the external API",
				Err:        errCustomer,
			}
		}

		WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
		var workOrderResponses transactionworkshoppayloads.WorkOrderResponse
		errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
		if errWorkOrder != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order data from the external API",
				Err:        errWorkOrder,
			}
		}

		convertedResponses = append(convertedResponses, transactionbodyshoppayloads.QualityControlResponse{
			WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
			WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
			VehicleCode:             vehicleResponses.Data.Master.VehicleChassisNumber,
			VehicleTnkb:             vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
			CustomerName:            customerResponses.CustomerName,
			WorkOrderSystemNumber:   entity.WorkOrderSystemNumber,
			VarianCode:              variantResponses.VariantCode,
			ModelCode:               modelResponses.ModelCode,
		})
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":            response.WorkOrderDocumentNumber,
			"work_order_date":                       response.WorkOrderDate,
			"model_code":                            response.ModelCode,
			"varian_code":                           response.VarianCode,
			"vehicle_chassis_number":                response.VehicleCode,
			"vehicle_registration_certificate_tnkb": response.VehicleTnkb,
			"customer_name":                         response.CustomerName,
			"work_order_system_number":              response.WorkOrderSystemNumber,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	pages.Rows = paginatedData
	pages.TotalRows = int64(totalRows)
	pages.TotalPages = totalPages

	return pages, nil
}

func (r *QualityControlBodyshopRepositoryImpl) GetById(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {
	var entity transactionbodyshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, entity)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_system_number = ? AND work_order_status_id IN (?,?)", id, utils.WoStatStop, utils.WoStatOngoing)

	if err := whereQuery.Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	// Fetch data model from external API
	modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external API
	variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data vehicle from external API
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Fetch data customer from external API
	CustomerUrl := config.EnvConfigs.SalesServiceUrl + "customer/" + strconv.Itoa(entity.CustomerId)
	var customerResponses transactionworkshoppayloads.CustomerResponse
	errCustomer := utils.Get(CustomerUrl, &customerResponses, nil)
	if errCustomer != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve customer data from the external API",
			Err:        errCustomer,
		}
	}

	// Fetch data work order from external API
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(entity.WorkOrderSystemNumber)
	var workOrderResponses transactionworkshoppayloads.WorkOrderResponse
	errWorkOrder := utils.Get(WorkOrderUrl, &workOrderResponses, nil)
	if errWorkOrder != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	var qualitycontrolDetails []transactionbodyshoppayloads.QualityControlDetailResponse
	var totalRows int64
	totalRowsQuery := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", entity.WorkOrderSystemNumber).
		Count(&totalRows).Error
	if totalRowsQuery != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count quality control details",
			Err:        totalRowsQuery,
		}
	}

	// Fetch paginated qc details
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Joins("INNER JOIN mtr_item AS WTA ON trx_work_order_detail.item_id = WTA.item_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_service_status AS MSS ON trx_work_order_detail.service_status_id = MSS.service_status_id").
		Select("WTA.item_code as operation_item_code, WTA.item_name as operation_item_name, trx_work_order_detail.frt_quantity as frt, trx_work_order_detail.service_status_id, MSS.service_status_description as service_status_name").
		Where("trx_work_order_detail.work_order_system_number = ?", id).
		Offset(pages.GetOffset()).
		Limit(pages.GetLimit())

	if err := query.Find(&qualitycontrolDetails).Error; err != nil {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service details",
			Err:        err,
		}
	}

	// Check if the service_status_id is valid
	if len(qualitycontrolDetails) == 0 {
		return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Quality control details not found",
		}
	}

	validStatuses := map[int]bool{
		utils.SrvStatStop:    true,
		utils.SrvStatQcPass:  true,
		utils.SrvStatReOrder: true,
	}
	for _, detail := range qualitycontrolDetails {
		if !validStatuses[detail.ServiceStatusId] {
			return transactionbodyshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
			}
		}
	}

	response := transactionbodyshoppayloads.QualityControlIdResponse{
		WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
		WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
		ModelName:               modelResponses.ModelName,
		VariantDescription:      variantResponses.VariantDescription,
		VehicleCode:             vehicleResponses.Data.Master.VehicleChassisNumber,
		VehicleTnkb:             vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
		CustomerName:            customerResponses.CustomerName,
		QualityControlDetails: transactionbodyshoppayloads.QualityControlDetailsResponse{
			Page:       pages.GetPage(),
			Limit:      pages.GetLimit(),
			TotalPages: int(math.Ceil(float64(totalRows) / float64(pages.GetLimit()))),
			TotalRows:  int(totalRows),
			Data:       qualitycontrolDetails,
		},
	}

	return response, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 2
// USE IN MODUL : AWS - 006  UPDATE DATA BY KEY (QC PASS) - GENERAL REPAIR
func (r *QualityControlBodyshopRepositoryImpl) Qcpass(tx *gorm.DB, id int, iddet int) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
	// Define variables
	var (
		currentStatus     int
		techAllocSysNo    int
		lineTypeOperation = 1
	)

	// Check the current WO_OPR_STATUS
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("work_order_status_id").
		Where("work_order_system_number = ? AND work_order_operation_item_line = ?", id, lineTypeOperation).
		First(&currentStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	// Validate the status
	if currentStatus != utils.WoStatStop {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "The current status of the work order is not valid",
		}
	}

	// Fetch work order details
	var details struct {
		VehicleId   int    `gorm:"column:vehicle_id"`
		BrandId     int    `gorm:"column:brand_id"`
		CompanyId   int    `gorm:"column:company_id"`
		OprItemCode string `gorm:"column:operation_item_code"`
		WoStatus    int    `gorm:"column:work_order_status_id"`
	}

	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("trx_work_order.vehicle_id, trx_work_order.company_id, trx_work_order_detail.operation_item_code, trx_work_order.work_order_status_id").
		Joins("JOIN trx_work_order ON trx_work_order_detail.work_order_system_number = trx_work_order.work_order_system_number").
		Where("trx_work_order_detail.work_order_system_number = ? AND trx_work_order_detail.work_order_operation_item_line = ?", id, lineTypeOperation).
		Scan(&details).Error
	if err != nil {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order details",
			Err:        err,
		}
	}

	// Fetch vehicle master data
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(details.VehicleId)
	if vehicleErr != nil {
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Fetch the latest TechAllocSystemNumber
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("ISNULL(MAX(technician_allocation_system_number), 0)").
		Where("work_order_system_number = ?", id).
		Where("work_order_line = ?", lineTypeOperation).
		Where("brand_id = ?", vehicleResponses.Data.Master.VehicleBrandID).
		Where("company_id = ?", details.CompanyId).
		Where("operation_code = ?", details.OprItemCode).
		Scan(&techAllocSysNo).Error
	if err != nil {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch the latest TechAllocSystemNumber",
			Err:        err,
		}
	}

	// Update WorkOrderDetail
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? and work_order_detail_id = ?", id, lineTypeOperation, iddet).
		Update("service_status_id", utils.WoStatQC).
		Error
	if err != nil {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update WorkOrderDetail",
			Err:        err,
		}
	}

	// Update WorkOrderAllocation
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("work_order_system_number = ? AND work_order_line = ?", id, lineTypeOperation).
		Update("service_status_id", utils.SrvStatQcPass).
		Error
	if err != nil {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update WorkOrderAllocation",
			Err:        err,
		}
	}

	// Check if all related items are updated
	var statusCount int64
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND service_status_id != ?", id, utils.SrvStatQcPass).
		Count(&statusCount).Error
	if err != nil {

		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count non-QC pass items",
			Err:        err,
		}
	}

	if statusCount == 0 {
		// Update WorkOrder if all related WorkOrderDetail have service_status_id as utils.SrvStatQcPass
		err = tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Update("work_order_status_id", utils.WoStatQC).
			Error
		if err != nil {

			return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update WorkOrder",
				Err:        err,
			}
		}

		// Update WorkOrder Detail
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ?", id).
			Update("work_order_status_id", utils.WoStatQC).
			Error
		if err != nil {

			return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update WorkOrder",
				Err:        err,
			}
		}
	}

	// Return response
	response := transactionbodyshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatQcPass,
		WorkOrderStatusName:   "QC Passed",
	}

	return response, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 4
// USE IN MODUL : AWS-006 SHEET: UPDATE DATA BY KEY (REORDER) - BODY REPAIR
func (r *QualityControlBodyshopRepositoryImpl) Reorder(tx *gorm.DB, id int, iddet int, payload transactionbodyshoppayloads.QualityControlReorder) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
	var (
		lineTypeOperation = 1
		woLine            = 1
	)

	// Check if the current WO_OPR_STATUS is valid
	var currentStatus int
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("service_status_id").
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND work_order_detail_id = ?", id, lineTypeOperation, iddet). // Assuming 1 is the value for Wo_Opr_Item_Line
		First(&currentStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	if currentStatus != utils.SrvStatStop {
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Operation Status is not valid, There are other technicians that is not Stop",
		}
	}

	// Update atWOTechAlloc
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("work_order_system_number = ? AND work_order_line = ?", id, woLine). // Assuming 1 is the value for Wo_Opr_Item_Line
		Updates(map[string]interface{}{
			"re_order": true,
		}).Error
	if err != nil {
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update atWOTechAlloc",
			Err:        err,
		}
	}

	// Update wtWorkOrder2
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_operation_item_line = ? and work_order_detail_id = ?", id, lineTypeOperation, iddet). // Assuming 1 is the value for Wo_Opr_Item_Line
		Updates(map[string]interface{}{
			"service_status_id":               utils.SrvStatReOrder,
			"quality_control_extra_frt":       payload.ExtraTime,
			"quality_control_total_extra_frt": gorm.Expr("quality_control_total_extra_frt + ?", payload.ExtraTime),
			"quality_control_extra_reason":    payload.Reason,
			"reorder_number":                  gorm.Expr("ISNULL(reorder_number, 0) + 1"),
		}).Error
	if err != nil {
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update wtWorkOrder2",
			Err:        err,
		}
	}

	// Update wtWorkOrder0
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", id).
		Updates(map[string]interface{}{
			"work_order_status_id": utils.WoStatOngoing,
		}).Error
	if err != nil {
		return transactionbodyshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update wtWorkOrder0",
			Err:        err,
		}
	}

	//return a response
	response := transactionbodyshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatReOrder,
		WorkOrderStatusName:   "ReOrder",
	}

	return response, nil
}
