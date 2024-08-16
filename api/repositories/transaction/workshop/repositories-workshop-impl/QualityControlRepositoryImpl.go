package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type QualityControlRepositoryImpl struct {
}

func OpenQualityControlRepositoryImpl() transactionworkshoprepository.QualityControlRepository {
	return &QualityControlRepositoryImpl{}
}

// uspg_wtWorkOrder0_Select
// IF @Option = 7
// USE IN MODUL : AWS - 006 QUALITY CONTROL PAGE 1 REQ: ???
func (r *QualityControlRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var entities []transactionworkshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.QualityControlRequest{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_status_id = ?", utils.WoStatStop) // 40 Stop

	if err := whereQuery.Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No entities found",
		}
	}

	var convertedResponses []transactionworkshoppayloads.QualityControlResponse

	for _, entity := range entities {
		// Fetch data brand from external API
		BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
		var brandResponses transactionworkshoppayloads.WorkOrderVehicleBrand
		errBrand := utils.Get(BrandUrl, &brandResponses, nil)
		if errBrand != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve brand data from the external API",
				Err:        errBrand,
			}
		}

		// Fetch data model from external API
		ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
		var modelResponses transactionworkshoppayloads.WorkOrderVehicleModel
		errModel := utils.Get(ModelUrl, &modelResponses, nil)
		if errModel != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve model data from the external API",
				Err:        errModel,
			}
		}

		// Fetch data variant from external API
		VariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entity.VariantId)
		var variantResponses transactionworkshoppayloads.WorkOrderVehicleVariant
		errVariant := utils.Get(VariantUrl, &variantResponses, nil)
		if errVariant != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve variant data from the external API",
				Err:        errVariant,
			}
		}

		// Fetch data vehicle from external API
		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
		var vehicleResponses transactionworkshoppayloads.VehicleResponse
		errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
		if errVehicle != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        errVehicle,
			}
		}

		// Fetch data customer from external API
		CustomerUrl := config.EnvConfigs.SalesServiceUrl + "customer/" + strconv.Itoa(entity.CustomerId)
		var customerResponses transactionworkshoppayloads.CustomerResponse
		errCustomer := utils.Get(CustomerUrl, &customerResponses, nil)
		if errCustomer != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
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
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order data from the external API",
				Err:        errWorkOrder,
			}
		}

		// Append converted response
		convertedResponses = append(convertedResponses, transactionworkshoppayloads.QualityControlResponse{
			WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
			WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
			ModelName:               modelResponses.ModelName,
			VariantName:             variantResponses.VariantName,
			VehicleCode:             vehicleResponses.VehicleCode,
			VehicleTnkb:             vehicleResponses.VehicleTnkb,
			CustomerName:            customerResponses.CustomerName,
		})
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":            response.WorkOrderDocumentNumber,
			"work_order_date":                       response.WorkOrderDate,
			"model_name":                            response.ModelName,
			"variant_name":                          response.VariantName,
			"vehicle_chassis_number":                response.VehicleCode,
			"vehicle_registration_certificate_tnkb": response.VehicleTnkb,
			"customer_name":                         response.CustomerName,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *QualityControlRepositoryImpl) GetById(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshoppayloads.QualityControlRequest

	joinTable := utils.CreateJoinSelectStatement(tx, entity)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("work_order_system_number = ? AND work_order_status_id IN (?,?)", id, utils.WoStatStop, utils.WoStatOngoing)

	if err := whereQuery.Find(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch entity",
			Err:        err,
		}
	}

	// Fetch data brand from external API
	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponses transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(BrandUrl, &brandResponses, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        errBrand,
		}
	}

	// Fetch data model from external API
	ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
	var modelResponses transactionworkshoppayloads.WorkOrderVehicleModel
	errModel := utils.Get(ModelUrl, &modelResponses, nil)
	if errModel != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        errModel,
		}
	}

	// Fetch data variant from external API
	VariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entity.VariantId)
	var variantResponses transactionworkshoppayloads.WorkOrderVehicleVariant
	errVariant := utils.Get(VariantUrl, &variantResponses, nil)
	if errVariant != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// Fetch data vehicle from external API
	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch data customer from external API
	CustomerUrl := config.EnvConfigs.SalesServiceUrl + "customer/" + strconv.Itoa(entity.CustomerId)
	var customerResponses transactionworkshoppayloads.CustomerResponse
	errCustomer := utils.Get(CustomerUrl, &customerResponses, nil)
	if errCustomer != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        errWorkOrder,
		}
	}

	var qualitycontrolDetails []transactionworkshoppayloads.QualityControlDetailResponse
	var totalRows int64
	totalRowsQuery := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", entity.WorkOrderSystemNumber).
		Count(&totalRows).Error
	if totalRowsQuery != nil {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get service details",
			Err:        err,
		}
	}

	// Check if the service_status_id is valid
	if len(qualitycontrolDetails) == 0 {
		return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
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
			return transactionworkshoppayloads.QualityControlIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
			}
		}
	}

	response := transactionworkshoppayloads.QualityControlIdResponse{
		WorkOrderDocumentNumber: workOrderResponses.WorkOrderDocumentNumber,
		WorkOrderDate:           workOrderResponses.WorkOrderDate.Format(time.RFC3339),
		ModelName:               modelResponses.ModelName,
		VariantName:             variantResponses.VariantName,
		VehicleCode:             vehicleResponses.VehicleCode,
		VehicleTnkb:             vehicleResponses.VehicleTnkb,
		CustomerName:            customerResponses.CustomerName,
		QualityControlDetails: transactionworkshoppayloads.QualityControlDetailsResponse{
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
func (r *QualityControlRepositoryImpl) Qcpass(tx *gorm.DB, id int, iddet int) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
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

			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	// Validate the status
	if currentStatus != utils.WoStatStop {

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order details",
			Err:        err,
		}
	}

	// Fetch vehicle master data
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(details.VehicleId)
	var vehicleResponses transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.Get(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch the latest TechAllocSystemNumber
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Select("ISNULL(MAX(technician_allocation_system_number), 0)").
		Where("work_order_system_number = ?", id).
		Where("work_order_line = ?", lineTypeOperation).
		Where("brand_id = ?", vehicleResponses.VehicleBrandId).
		Where("company_id = ?", details.CompanyId).
		Where("operation_code = ?", details.OprItemCode).
		Scan(&techAllocSysNo).Error
	if err != nil {

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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

			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update WorkOrder",
				Err:        err,
			}
		}
	}

	// Return response
	response := transactionworkshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatQcPass,
		WorkOrderStatusName:   "QC Passed",
	}

	return response, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 1
// USE IN MODUL : AWS-006 SHEET: RE-ORDER
func (r *QualityControlRepositoryImpl) Reorder(tx *gorm.DB, id int, iddet int, payload transactionworkshoppayloads.QualityControlReorder) (transactionworkshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse) {
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
			return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Operation Status is not valid",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch operation status",
			Err:        err,
		}
	}

	if currentStatus != utils.SrvStatStop {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Operation Status is not valid",
		}
	}

	// Update atWOTechAlloc
	err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("work_order_system_number = ? AND work_order_line = ?", id, woLine). // Assuming 1 is the value for Wo_Opr_Item_Line
		Updates(map[string]interface{}{
			"re_order": true,
		}).Error
	if err != nil {
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.QualityControlUpdateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update wtWorkOrder0",
			Err:        err,
		}
	}

	//return a response
	response := transactionworkshoppayloads.QualityControlUpdateResponse{
		WorkOrderSystemNumber: id,
		WorkOrderDetailId:     iddet,
		WorkOrderStatusId:     utils.SrvStatReOrder,
		WorkOrderStatusName:   "ReOrder",
	}

	return response, nil
}
