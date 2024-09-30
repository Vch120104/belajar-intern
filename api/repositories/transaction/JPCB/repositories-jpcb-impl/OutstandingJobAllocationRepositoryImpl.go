package transactionjpcbrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteroperationentities "after-sales/api/entities/master/operation"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	"after-sales/api/utils"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OutstandingJobAllocationRepositoryImpl struct {
}

func StartOutStandingJobAllocationRepository() transactionjpcbrepository.OutstandingJobAllocationRepository {
	return &OutstandingJobAllocationRepositoryImpl{}
}

func (r *OutstandingJobAllocationRepositoryImpl) GetAllOutstandingJobAllocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entitiesBooking := transactionworkshopentities.BookingEstimation{}
	entitiesWorkOrder := transactionworkshopentities.WorkOrder{}
	responses := []transactionjpcbpayloads.OutstandingJobAllocationGetAllPayload{}

	bookingServiceDate := ""
	companyId := 0
	var referenceDocumentType *string
	var referenceDocumentNumber *string
	var plateNumber *string
	for _, filter := range filterCondition {
		if filter.ColumnField == "booking_service_date" {
			bookingServiceDate = filter.ColumnValue
		}
		if filter.ColumnField == "company_id" {
			companyId, _ = strconv.Atoi(filter.ColumnValue)
		}
		if filter.ColumnField == "reference_document_type" && filter.ColumnValue != "" {
			referenceDocumentType = &filter.ColumnValue
		}
		if filter.ColumnField == "reference_document_number" && filter.ColumnValue != "" {
			referenceDocumentNumber = &filter.ColumnValue
		}
		if filter.ColumnField == "tnkb" && filter.ColumnValue != "" {
			plateNumber = &filter.ColumnValue
		}
	}

	profitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center-by-name/Workshop"
	profitcenterPayloads := transactionjpcbpayloads.OutstandingJAProfitCenterPayload{}
	if err := utils.Get(profitCenterUrl, &profitcenterPayloads, nil); err != nil || profitcenterPayloads.ProfitCenterId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching profit center data"),
		}
	}

	approvalStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status-description/Wait%20Approve"
	approvalStatusPayloads := []transactionjpcbpayloads.OutstandingJAApprovalStatusPayload{}
	if err := utils.GetArray(approvalStatusUrl, &approvalStatusPayloads, nil); err != nil || approvalStatusPayloads[0].ApprovalStatusId == "" {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching approval status data"),
		}
	}
	approvalStatusId, _ := strconv.Atoi(approvalStatusPayloads[0].ApprovalStatusId)

	documentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-description/New%20Document"
	documentStatusPayloads := transactionjpcbpayloads.OutstandingJADocumentStatusPayload{}
	if err := utils.Get(documentStatusUrl, &documentStatusPayloads, nil); err != nil || documentStatusPayloads.DocumentStatusId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching document status data"),
		}
	}

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type-by-name/Operation"
	lineTypePayloads := transactionjpcbpayloads.OutstandingJALineTypePayload{}
	if err := utils.Get(lineTypeUrl, &lineTypePayloads, nil); err != nil || lineTypePayloads.LineTypeId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching line type data"),
		}
	}

	workOrderStatusUrl := config.EnvConfigs.GeneralServiceUrl + "work-order-status-by-descriptions/Draft,QC Pass,Closed"
	workOrderStatusPayloads := []transactionjpcbpayloads.OutstandingJAWorkOrderStatusPayload{}
	if err := utils.GetArray(workOrderStatusUrl, &workOrderStatusPayloads, nil); err != nil || len(workOrderStatusPayloads) != 3 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching work order status data"),
		}
	}
	workOrderStatusIds := []int{}
	for _, workOrder := range workOrderStatusPayloads {
		workOrderStatusIds = append(workOrderStatusIds, workOrder.WorkOrderStatusId)
	}

	serviceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status-by-description/Transfer"
	serviceStatusPayloads := transactionjpcbpayloads.OutstandingJAServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayloads, nil); err != nil || serviceStatusPayloads.ServiceStatusId == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status data"),
		}
	}

	baseModelQuery1 := tx.Model(&entitiesBooking).
		Select(`
			'BOOKING' AS reference_document_type,
			trx_booking_estimation.batch_system_number AS reference_system_number,
			tbea.booking_document_number AS reference_document_number,
			trx_booking_estimation.vehicle_id AS vehicle_id,
			'' leave_car,
			'' is_from_booking,
			moc.operation_code + '-' + moc.operation_name AS operation_description,
			tbed.frt_quantity AS reference_frt,
			tbea.booking_service_date AS promise_time,
			0 last_release_by,
			'' remark`).
		Joins("LEFT JOIN trx_booking_estimation_allocation tbea ON tbea.booking_system_number = trx_booking_estimation.booking_system_number").
		Joins("LEFT JOIN trx_booking_estimation_detail tbed ON tbed.estimation_system_number = trx_booking_estimation.estimation_system_number").
		Joins("LEFT JOIN mtr_operation_code moc ON moc.operation_id = tbed.item_operation_id").
		Joins("LEFT JOIN trx_work_order_allocation twoa ON twoa.booking_system_number = trx_booking_estimation.booking_system_number").
		Joins("LEFT JOIN trx_booking_estimation_service_discount tbesd ON tbesd.estimation_system_number = trx_booking_estimation.estimation_system_number").
		Where("technician_allocation_system_number IS NULL OR technician_allocation_system_number = 0").
		Where("trx_booking_estimation.booking_system_number != 0 AND trx_booking_estimation.booking_system_number IS NOT NULL").
		Where("trx_booking_estimation.profit_center_id = ?", profitcenterPayloads.ProfitCenterId).
		Where("tbesd.estimation_discount_approval_status != ?", approvalStatusId).
		Where("(CAST(tbea.booking_service_date AS DATE) LIKE '?' AND tbea.document_status_id = ?) OR (ISNULL(tbea.booking_service_date, '') = '' AND (trx_booking_estimation.booking_system_number > 0 AND tbesd.estimation_discount_approval_status = ?))", bookingServiceDate, documentStatusPayloads.DocumentStatusId, documentStatusPayloads.DocumentStatusId).
		Where("tbed.line_type_id = ?", lineTypePayloads.LineTypeId).
		Where("tbea.booking_document_number != ''")

	baseModelQuery2 := tx.Model(&entitiesWorkOrder).
		Select(`
			'WORKORDER' AS reference_document_type,
			trx_work_order.work_order_system_number AS reference_system_number,
			trx_work_order.work_order_document_number AS reference_document_number,
			trx_work_order.vehicle_id AS vehicle_id,
			CASE WHEN trx_work_order.leave_car = 1
				THEN 'YES'
				ELSE 'NO'
			END leave_car,
			CASE WHEN trx_work_order.booking_system_number > 0
				THEN 'YES'
				ELSE 'NO'
			END is_from_booking,
			twod.operation_item_code AS operation_description,
			CASE WHEN mic.item_class_code = 'OJ'
				THEN 0
				ELSE twod.frt_quantity
			END reference_frt,
			FORMAT(trx_work_order.promise_date, 'd MMMM yyyy') + ' : ' + FORMAT(trx_work_order.promise_time, 'HH:mm') AS promise_time,
			tsl.technician_id AS last_release_by,
			trx_work_order.remark AS remark
		`).
		Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = trx_work_order.work_order_system_number").
		Joins("LEFT JOIN mtr_item mi ON mi.item_id = twod.operation_item_id").
		Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id AND mic.item_class_code = 'OJ'").
		Joins("LEFT JOIN trx_service_log tsl ON tsl.work_order_system_number = trx_work_order.work_order_system_number AND tsl.operation_item_code = twod.operation_item_code AND tsl.service_status_id = ?", serviceStatusPayloads.ServiceStatusId).
		Joins("LEFT JOIN trx_work_order_allocation twoa ON twoa.work_order_system_number = trx_work_order.work_order_system_number AND twoa.operation_code = twod.operation_item_code AND ISNULL(twoa.re_order, 0) = 0").
		Where("trx_work_order.work_order_status_id NOT IN(?)", workOrderStatusIds).
		Where("ISNULL(twod.service_status_id, 0) IN (0, ?)", serviceStatusPayloads.ServiceStatusId).
		Where("twod.line_type_id = ? OR mic.item_class_code = 'OJ'", lineTypePayloads.LineTypeId).
		Where("twoa.technician_allocation_system_number IS NULL").
		Where("trx_work_order.company_id = ?", companyId)

	techAllocSysNum := tx.Table("trx_work_order_allocation AS alcc").
		Select("TOP 1 alcc.technician_allocation_system_number").
		Where("alc.work_order_system_number = alcc.work_order_system_number").
		Where("alc.work_order_line = alcc.work_order_line").
		Order("alcc.technician_allocation_system_number DESC")

	techAlloc := tx.Table("trx_work_order_allocation AS alc").
		Where("ISNULL(alc.re_order, 0) = 0").
		Where("alc.work_order_system_number = trx_work_order.work_order_system_number").
		Where("alc.work_order_line = twod.work_order_operation_item_line").
		Where("alc.technician_allocation_system_number = (?)", techAllocSysNum)

	baseModelQuery3 := tx.Model(&entitiesWorkOrder).
		Select(`
			'WORKORDER' AS reference_document_type,
			trx_work_order.work_order_system_number AS reference_system_number,
			trx_work_order.work_order_document_number AS reference_document_number,
			trx_work_order.vehicle_id AS vehicle_id,
			CASE WHEN trx_work_order.leave_car = 1
				THEN 'YES'
				ELSE 'NO'
			END leave_car,
			CASE WHEN trx_work_order.booking_system_number > 0
				THEN 'YES'
				ELSE 'NO'
			END is_from_booking,
			twod.operation_item_code AS operation_description,
			CASE WHEN mic.item_class_code = 'OJ'
				THEN 0
				ELSE twod.frt_quantity
			END reference_frt,
			FORMAT(trx_work_order.promise_date, 'd MMMM yyyy') + ' : ' + FORMAT(trx_work_order.promise_time, 'HH:mm') AS promise_time,
			tsl.technician_id AS last_release_by,
			trx_work_order.remark AS remark
		`).
		Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = trx_work_order.work_order_system_number").
		Joins("LEFT JOIN mtr_item mi ON mi.item_id = twod.operation_item_id").
		Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id AND mic.item_class_code = 'OJ'").
		Joins("LEFT JOIN trx_service_log tsl ON tsl.work_order_system_number = trx_work_order.work_order_system_number AND tsl.operation_item_code = twod.operation_item_code AND tsl.service_status_id = ?", serviceStatusPayloads.ServiceStatusId).
		Joins("LEFT JOIN trx_work_order_allocation twoa ON twoa.work_order_system_number = trx_work_order.work_order_system_number AND twoa.operation_code = twod.operation_item_code AND ISNULL(twoa.re_order, 0) = ISNULL(twod.reorder_number, 0) AND ISNULL(twoa.re_order, 0) = 0").
		Where("trx_work_order.work_order_status_id NOT IN(?)", workOrderStatusIds).
		Where("twod.line_type_id = ?", lineTypePayloads.LineTypeId).
		Where("ISNULL(twod.reorder_number, 0) > 0").
		Where("trx_work_order.company_id = ?", companyId).
		Where("NOT EXISTS (?)", techAlloc)

	addCondition := "(reference_document_type LIKE 'BOOKING' OR reference_document_type LIKE 'WORKORDER')"
	values := []interface{}{baseModelQuery1, baseModelQuery2, baseModelQuery3}
	if referenceDocumentType != nil {
		addCondition += " AND reference_document_type = ?"
		values = append(values, referenceDocumentType)
	}
	if referenceDocumentNumber != nil {
		addCondition += " AND reference_document_number = ?"
		values = append(values, referenceDocumentNumber)
	}
	if plateNumber != nil {
		vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_registration_certificate_tnkb=" + *plateNumber
		vehiclePayloads := []transactionjpcbpayloads.OutstandingJAVehicleMasterByTnkbPayload{}
		if err := utils.GetArray(vehicleUrl, &vehiclePayloads, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("error fetching vehicle by tnkb data"),
			}
		}
		if len(vehiclePayloads) > 0 {
			vehicleIds := []int{}
			for _, vehicle := range vehiclePayloads {
				vehicleIds = append(vehicleIds, vehicle.VehicleId)
			}
			addCondition += " AND vehicle_id IN (?)"
			values = append(values, vehicleIds)
		} else {
			addCondition += " AND vehicle_id = ?"
			values = append(values, -1)
		}
	}
	values = append(values, pages.GetOffset(), pages.GetLimit())

	err := tx.Raw(`SELECT * FROM (? UNION ALL ? UNION ALL ?) AS SRC WHERE `+addCondition+` ORDER BY promise_time OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`, values...).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var totalRows int64
	err = tx.Raw(`SELECT COUNT(*) AS total_rows FROM (? UNION ALL ? UNION ALL ?) AS SRC`, baseModelQuery1, baseModelQuery2, baseModelQuery3).Pluck("total_rows", &totalRows).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	finalResponse := []transactionjpcbpayloads.OutstandingJobAllocationGetAllResponse{}
	for _, data := range responses {
		Tnkb := ""
		if *data.VehicleId != 0 && data.VehicleId != nil {
			vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(*data.VehicleId)
			vehiclePayloads := transactionjpcbpayloads.OutstandingJAVehicleByIdPayload{}
			if err := utils.Get(vehicleUrl, &vehiclePayloads, nil); err != nil || vehiclePayloads.Master.VehicleId == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("error fetching vehicle data"),
				}
			}
			Tnkb = vehiclePayloads.Stnk.VehicleRegistrationCertificateTnkb
		}

		employeeName := ""
		if *data.LastReleaseBy != 0 && data.LastReleaseBy != nil {
			employeeUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(*data.LastReleaseBy)
			employeePayloads := transactionjpcbpayloads.OutstandingJAEmployeePayload{}
			if err := utils.Get(employeeUrl, &employeePayloads, nil); err != nil || employeePayloads.UserEmployeeId == 0 {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("error fetching employee data"),
				}
			}
			employeeName = employeePayloads.EmployeeName
		}

		result := transactionjpcbpayloads.OutstandingJobAllocationGetAllResponse{
			ReferenceDocumentType:   data.ReferenceDocumentType,
			ReferenceSystemNumber:   data.ReferenceSystemNumber,
			ReferenceDocumentNumber: data.ReferenceDocumentNumber,
			TNKB:                    &Tnkb,
			LeaveCar:                data.LeaveCar,
			IsFromBooking:           data.IsFromBooking,
			OperationDescription:    data.OperationDescription,
			ReferenceFRT:            data.ReferenceFRT,
			PromiseTime:             data.PromiseTime,
			LastReleaseBy:           &employeeName,
			Remark:                  data.Remark,
		}
		finalResponse = append(finalResponse, result)
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))
	pages.Rows = finalResponse

	return pages, nil
}

func (r *OutstandingJobAllocationRepositoryImpl) GetByTypeIdOutstandingJobAllocation(tx *gorm.DB, referenceDocumentType string, referenceSystemNumber int) (transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse, *exceptions.BaseErrorResponse) {
	entitiesBooking := transactionworkshopentities.BookingEstimation{}
	entitiesWorkOrder := transactionworkshopentities.WorkOrder{}
	response := transactionjpcbpayloads.OutstandingJobAllocationGetByTypeIdResponse{}

	var err error
	if referenceDocumentType == "BOOKING" {
		err = tx.Model(&entitiesBooking).
			Select(`
				trx_booking_estimation.batch_system_number AS reference_system_number,
				tbea.booking_document_number AS reference_document_number,
				moc.operation_code + '-' + moc.operation_name AS operation_description,
				tbed.frt_quantity AS reference_frt,
				0 AS job_progress
			`).
			Joins("INNER JOIN trx_booking_estimation_allocation tbea ON tbea.booking_system_number = trx_booking_estimation.booking_system_number").
			Joins("LEFT JOIN trx_booking_estimation_detail tbed ON tbed.estimation_system_number = trx_booking_estimation.estimation_system_number").
			Joins("LEFT JOIN mtr_operation_code moc ON moc.operation_id = tbed.item_operation_id").
			Where("trx_booking_estimation.batch_system_number = ?", referenceSystemNumber).
			First(&response).Error
	} else if referenceDocumentType == "WORKORDER" {
		err = tx.Model(&entitiesWorkOrder).
			Select(`
				trx_work_order.work_order_system_number AS reference_system_number,
				trx_work_order.work_order_document_number AS reference_document_number,
				twod.operation_item_code AS operation_description,
				CASE WHEN mic.item_class_code = 'OJ'
					THEN 0
					ELSE twod.frt_quantity
				END reference_frt,
				0 AS job_progress
			`).
			Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = trx_work_order.work_order_system_number").
			Joins("LEFT JOIN mtr_item mi ON mi.item_id = twod.operation_item_id").
			Joins("INNER JOIN mtr_item_class mic ON mic.item_class_id = mi.item_class_id AND mic.item_class_code = 'OJ'").
			Where("trx_work_order.work_order_system_number = ?", referenceSystemNumber).
			First(&response).Error
	} else {
		err = errors.New("type must be either 'booking' or 'workorder'")
	}

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

// uspg_atWoTechAlloc_Insert
// IF @Option = 2
func (r *OutstandingJobAllocationRepositoryImpl) SaveOutstandingJobAllocation(tx *gorm.DB, referenceDocumentType string, referenceSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationSaveRequest, operationCodeResult masteroperationpayloads.OperationCodeResponse) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest, *exceptions.BaseErrorResponse) {
	entitiesWorkOrderAllocation := transactionworkshopentities.WorkOrderAllocation{}
	entitiesWorkOrder := transactionworkshopentities.WorkOrder{}
	entitiesBookingEstimation := transactionworkshopentities.BookingEstimation{}
	entitiesOperationModelMapping := masteroperationentities.OperationModelMapping{}
	response := transactionjpcbpayloads.SettingTechnicianGetByIdResponse{}
	responseUpdate := transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest{}

	factorX := 1.0
	employeeUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(req.UserEmployeeId)
	employeePayloads := transactionjpcbpayloads.OutstandingJAEmployeePayload{}
	if err := utils.Get(employeeUrl, &employeePayloads, nil); err != nil {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching employee data"),
		}
	}
	if employeePayloads.FactorX != 0 {
		factorX = employeePayloads.FactorX
	}

	serviceDate := req.ServiceDate
	techAllocStartDate := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, serviceDate.Location())

	oriSequenceNumber := 0
	err := tx.Model(&entitiesWorkOrderAllocation).
		Select("ISNULL(MAX(sequence_number), 0) + 1 AS sequence_number").
		Where("company_id = ?", req.CompanyId).
		Where("technician_id = ?", req.UserEmployeeId).
		Where("CAST(tech_alloc_start_date AS DATE) = ?", techAllocStartDate).
		Pluck("sequence_number", &oriSequenceNumber).Error
	if err != nil || oriSequenceNumber == 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching sequence number data"),
		}
	}

	var frtQuantityJPCB float64
	responseReferenceData := transactionjpcbpayloads.OutstandingJobAllocationReferencePayload{}
	if referenceDocumentType == "WORKORDER" {
		err = tx.Model(&entitiesWorkOrder).
			Select(`
				'WORKORDER' AS reference_document_type,
				trx_work_order.work_order_system_number AS reference_system_number,
				trx_work_order.work_order_document_number AS reference_document_number,
				trx_work_order.work_order_date AS reference_document_date,
				trx_work_order.cost_center_id AS cost_profit_center_id,
				trx_work_order.vehicle_id,
				twod.work_order_operation_item_line AS line,
				ISNULL(twod.frt_quantity, 0),
				trx_work_order.work_order_date,
				twod.quality_control_extra_frt,
				twod.reorder_number
			`).
			Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = trx_work_order.work_order_system_number").
			Where("trx_work_order.work_order_system_number = ?", referenceSystemNumber).
			First(&responseReferenceData).Error
		if err != nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if *responseReferenceData.ReorderNumber > 0 {
			if responseReferenceData.QualityControlExtraFrt == nil {
				return response, responseUpdate, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("quality control extra frt data is null"),
				}
			}
			frtQuantityJPCB = *responseReferenceData.QualityControlExtraFrt
		} else {
			if responseReferenceData.FrtQuantity == nil {
				return response, responseUpdate, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("frt quantity data is null"),
				}
			}
			frtQuantityJPCB = *responseReferenceData.FrtQuantity * factorX
		}
	} else if referenceDocumentType == "BOOKING" {
		err = tx.Model(&entitiesBookingEstimation).
			Select(`
				'BOOKING' AS reference_document_type,
				trx_booking_estimation.batch_system_number AS reference_system_number,
				tbea.booking_document_number AS reference_document_number,
				tbea.booking_date AS reference_document_number,
				trx_booking_estimation.profit_center_id AS cost_profit_center_id,
				trx_booking_estimation.vehicle_id,
				tbed.estimation_line_code AS line,
				tbed.frt_quantity,
				NULL AS work_order_date,
				tbea.booking_service_time
			`).
			Joins("LEFT JOIN trx_booking_estimation_allocation tbea ON tbea.booking_system_number = trx_booking_estimation.booking_system_number").
			Joins("LEFT JOIN trx_booking_estimation_detail tbed ON tbed.estimation_system_number = trx_booking_estimation.estimation_system_number").
			Where("trx_booking_estimation.batch_system_number = ?", referenceSystemNumber).
			First(&responseReferenceData).Error
		if err != nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		} else if responseReferenceData.BookingServiceTime == nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("booking service time is empty"),
			}
		} else if responseReferenceData.FrtQuantity == nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("frt quantity data is empty"),
			}
		}
		frtQuantityJPCB = *responseReferenceData.FrtQuantity * factorX
	} else {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("document type must be either 'workorder' or 'booking'"),
		}
	}

	if !(*responseReferenceData.VehicleId > 0) {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("referenced data doesn't have vehicle"),
		}
	}

	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(*responseReferenceData.VehicleId)
	vehiclePayload := transactionjpcbpayloads.OutstandingJAVehicleByIdPayload{}
	if err := utils.Get(vehicleUrl, &vehiclePayload, nil); err != nil || vehiclePayload.Master.VehicleId == 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching vehicle data"),
		}
	}

	if req.IsExpress {
		var frtHourExpress float64
		err = tx.Model(&entitiesOperationModelMapping).
			Select("mof.frt_hour_express").
			Joins("INNER JOIN mtr_operation_frt mof ON mof.operation_model_mapping_id = mtr_operation_model_mapping.operation_model_mapping_id").
			Where("mtr_operation_model_mapping.brand_id = ?", vehiclePayload.Master.VehicleBrandId).
			Where("mtr_operation_model_mapping.model_id = ?", vehiclePayload.Master.VehicleModelId).
			Where("mtr_operation_model_mapping.operation_id = ?", req.OperationId).
			Where("mof.variant_id = ?", vehiclePayload.Master.VehicleVariantId).
			First(&frtHourExpress).Error
		if err != nil || frtHourExpress == 0 {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("FRT Express for Operation code is not set yet, Please Setup FRT Express for this Operation code first"),
			}
		}
		frtQuantityJPCB = frtHourExpress
	}

	setTechDetailEntities := transactionjpcbentities.SettingTechnicianDetail{}
	var shiftGroupId int
	err = tx.Model(&setTechDetailEntities).
		Select("trx_setting_technician_detail.shift_group_id").
		Joins("INNER JOIN trx_setting_technician tst ON tst.setting_technician_system_number = trx_setting_technician_detail.setting_technician_system_number").
		Where("tst.company_id = ?", req.CompanyId).
		Where("trx_setting_technician_detail.technician_employee_number_id = ?", req.UserEmployeeId).
		Order("tst.effective_date DESC").
		First(&shiftGroupId).Error
	if err != nil || shiftGroupId == 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("there is no shift group"),
		}
	}

	var foremanId int
	employeeGroupUrl := config.EnvConfigs.GeneralServiceUrl + "employee-group-by-member-id/" + strconv.Itoa(req.CompanyId) + "/" + strconv.Itoa(req.UserEmployeeId)
	employeeGroupPayload := transactionjpcbpayloads.OutstandingJAEmployeeGroupPayload{}
	if err := utils.Get(employeeGroupUrl, &employeeGroupPayload, nil); err != nil || employeeGroupPayload.EmployeeGroupLeaderId == 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching foreman / employee group data"),
		}
	}
	foremanId = employeeGroupPayload.EmployeeGroupLeaderId

	techAllocEndDate := techAllocStartDate

	day := techAllocStartDate.Weekday()

	shiftScheduleEntities := masterentities.ShiftSchedule{}
	shiftScheduleQuery := tx.Model(&shiftScheduleEntities).Where(masterentities.ShiftSchedule{ShiftScheduleId: shiftGroupId})
	switch day {
	case 0:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Sunday: true})
	case 1:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Monday: true})
	case 2:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Tuesday: true})
	case 3:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Wednesday: true})
	case 4:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Thursday: true})
	case 5:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Friday: true})
	case 6:
		shiftScheduleQuery = shiftScheduleQuery.Where(masterentities.ShiftSchedule{Saturday: true})
	default:
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("something wrong with the tech alloc start date"),
		}
	}

	shiftScheduleQuery = shiftScheduleQuery.Order("effective_date DESC")
	err = shiftScheduleQuery.First(&shiftScheduleEntities).Error
	if err != nil {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	shiftCode := shiftScheduleEntities.ShiftCode
	techAllocStartTime := shiftScheduleEntities.StartTime
	restStartTime := shiftScheduleEntities.RestStartTime
	restEndTime := shiftScheduleEntities.RestEndTime

	if oriSequenceNumber > 1 {
		err = tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
			Select("tech_alloc_last_end_time").
			Where(transactionworkshopentities.WorkOrderAllocation{
				CompanyId:          req.CompanyId,
				TechnicianId:       req.UserEmployeeId,
				TechAllocStartDate: techAllocStartDate,
				SequenceNumber:     oriSequenceNumber - 1,
			}).
			First(&techAllocStartTime).Error
		if err != nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	currentDateTime := time.Now()
	currentDate := time.Date(currentDateTime.Year(), currentDateTime.Month(), currentDateTime.Day(), 0, 0, 0, 0, currentDateTime.Location())

	if oriSequenceNumber == 1 && techAllocStartDate == currentDate {
		timeNumeric := float64(currentDateTime.Hour()) + float64(currentDateTime.Minute())/100
		if techAllocStartTime < timeNumeric {
			techAllocStartTime = timeNumeric
		}
	}

	if referenceDocumentType == "BOOKING" {
		if techAllocStartTime < *responseReferenceData.BookingServiceTime {
			techAllocStartTime = *responseReferenceData.BookingServiceTime
		}
	}

	techAllocEndTime := techAllocStartTime + frtQuantityJPCB
	if (techAllocStartTime+frtQuantityJPCB) > restStartTime && techAllocStartTime < restEndTime {
		techAllocEndTime = restEndTime + ((techAllocStartTime + frtQuantityJPCB) - restStartTime)
	}

	if techAllocEndTime > 21.00 {
		techAllocEndTime = 21.00
	}

	serviceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status-by-description/Draft"
	serviceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayload, nil); err != nil || serviceStatusPayload.ServiceStatusId == 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status data"),
		}
	}

	var workOrderAllocationCount int64
	workOrderAllocationQuery := tx.Model(&transactionworkshopentities.WorkOrderAllocation{}).
		Where("company_id = ?", req.CompanyId).
		Where("profit_center_id = ?", *responseReferenceData.CostProfitCenterId).
		Where("operation_code = ?", operationCodeResult.OperationCode).
		Where("technician_id = ?", req.UserEmployeeId).
		Where("re_order = 0")
	if referenceDocumentType == "BOOKING" {
		workOrderAllocationQuery = workOrderAllocationQuery.Where("booking_system_number = ?", referenceSystemNumber)
	} else {
		workOrderAllocationQuery = workOrderAllocationQuery.Where("work_order_system_number = ?", referenceSystemNumber)
	}
	err = workOrderAllocationQuery.Count(&workOrderAllocationCount).Error
	if err != nil || workOrderAllocationCount > 0 {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("data already exist"),
		}
	}

	entitiesWorkOrderAllocation = transactionworkshopentities.WorkOrderAllocation{
		IsActive:                true,
		CompanyId:               req.CompanyId,
		BrandId:                 vehiclePayload.Master.VehicleBrandId,
		ProfitCenterId:          *responseReferenceData.CostProfitCenterId,
		TechnicianId:            req.UserEmployeeId,
		ForemanId:               foremanId,
		UsingGroup:              true,
		TechnicianGroupId:       employeeGroupPayload.EmployeeGroupMemberId,
		SequenceNumber:          oriSequenceNumber,
		TechAllocStartDate:      techAllocStartDate,
		TechAllocEndDate:        techAllocEndDate,
		TechAllocStartTime:      techAllocStartTime,
		TechAllocEndTime:        techAllocEndTime,
		TechAllocTotalTime:      frtQuantityJPCB,
		TechAllocLastStartDate:  techAllocStartDate,
		TechAllocLastEndDate:    techAllocEndDate,
		TechAllocLastStartTime:  techAllocStartTime,
		TechAllocLastEndTime:    techAllocEndTime,
		OperationCode:           operationCodeResult.OperationCode,
		ShiftCode:               shiftCode,
		ServActualTime:          0,
		ServPendingTime:         0,
		ServProgressTime:        0,
		ServTotalActualTime:     0,
		ServStatus:              serviceStatusPayload.ServiceStatusId,
		BookingSystemNumber:     0,
		BookingDocumentNumber:   "",
		BookingLine:             0,
		WorkOrderSystemNumber:   0,
		WorkOrderDocumentNumber: "",
		WorkOrderLine:           0,
		ReOrder:                 false,
		InvoiceSystemNumber:     0,
		InvoiceDocumentNumber:   "",
		IncentiveSystemNumber:   0,
		FactorX:                 factorX,
		IsExpress:               req.IsExpress,
		Frt:                     *responseReferenceData.FrtQuantity,
		BookingServiceTime:      0,
	}

	if referenceDocumentType == "BOOKING" {
		entitiesWorkOrderAllocation.BookingSystemNumber = *responseReferenceData.ReferenceSystemNumber
		entitiesWorkOrderAllocation.BookingDocumentNumber = ""
		if responseReferenceData.ReferenceDocumentNumber != nil {
			entitiesWorkOrderAllocation.BookingDocumentNumber = *responseReferenceData.ReferenceDocumentNumber
		}
		entitiesWorkOrderAllocation.BookingServiceTime = 0
		if responseReferenceData.BookingServiceTime != nil {
			entitiesWorkOrderAllocation.BookingServiceTime = *responseReferenceData.BookingServiceTime
		}
	} else {
		entitiesWorkOrderAllocation.WorkOrderSystemNumber = *responseReferenceData.ReferenceSystemNumber
		entitiesWorkOrderAllocation.WorkOrderDocumentNumber = *responseReferenceData.ReferenceDocumentNumber
	}

	err = tx.Save(&entitiesWorkOrderAllocation).Error
	if err != nil {
		return response, responseUpdate, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if referenceDocumentType == "WORKORDER" {
		entitiesServiceLog := transactionworkshopentities.ServiceLog{
			CompanyId:                        req.CompanyId,
			TechnicianAllocationSystemNumber: entitiesWorkOrderAllocation.TechAllocSystemNumber,
			WorkOrderSystemNumber:            entitiesWorkOrderAllocation.WorkOrderSystemNumber,
			WorkOrderDocumentNumber:          entitiesWorkOrderAllocation.WorkOrderDocumentNumber,
			WorkOrderOperationId:             req.OperationId,
			ShiftScheduleId:                  shiftGroupId,
			Frt:                              frtQuantityJPCB,
			ServiceStatusId:                  entitiesWorkOrderAllocation.ServStatus,
			StartDatetime:                    entitiesWorkOrderAllocation.TechAllocStartDate,
			EndDatetime:                      entitiesWorkOrderAllocation.TechAllocEndDate,
			ActualTime:                       0,
			EstimatedPendingTime:             0,
			PendingTime:                      0,
			SequenceNumber:                   entitiesWorkOrderAllocation.SequenceNumber,
		}
		err = tx.Save(&entitiesServiceLog).Error
		if err != nil {
			return response, responseUpdate, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	if oriSequenceNumber > 1 && oriSequenceNumber != req.SequenceNumber {
		responseUpdate = transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest{
			TechAllocSystemNumber: entitiesWorkOrderAllocation.TechAllocSystemNumber,
			CompanyId:             req.CompanyId,
			OriginalTechnicianId:  req.UserEmployeeId,
			TechnicianId:          req.UserEmployeeId,
			OriSequenceNumber:     oriSequenceNumber,
			SequenceNumber:        req.SequenceNumber,
			ServiceDate:           req.ServiceDate,
		}
	}

	return response, responseUpdate, nil
}

// uspg_uspg_atWoTechAlloc_Update
// IF @Option = 3
func (r *OutstandingJobAllocationRepositoryImpl) UpdateOutstandingJobAllocation(tx *gorm.DB, techAllocSystemNumber int, req transactionjpcbpayloads.OutstandingJobAllocationUpdateRequest) (transactionjpcbpayloads.OutstandingJobAllocationUpdateResponse, *exceptions.BaseErrorResponse) {
	response := transactionjpcbpayloads.OutstandingJobAllocationUpdateResponse{}

	// begin validate job status
	entitiesWorkOrderAllocation := transactionworkshopentities.WorkOrderAllocation{}
	serviceStatusList := []string{"Draft", "Transfer", "Auto Release"}
	sourceResponseWorkOrderAllocation := transactionjpcbpayloads.OutstandingJobAllocationSourceTargetPayload{}
	sourceCheckServiceStatus := false
	targetResponseWorkOrderAllocation := transactionjpcbpayloads.OutstandingJobAllocationSourceTargetPayload{}
	targetCheckServiceStatus := false

	err := tx.Model(&entitiesWorkOrderAllocation).
		Select(`
			trx_work_order_allocation.technician_allocation_system_number,
			trx_work_order_allocation.service_status_id,
			two.vehicle_id,
			trx_work_order_allocation.technician_id,
			trx_work_order_allocation.work_order_system_number,
			trx_work_order_allocation.work_order_document_number,
			trx_work_order_allocation.operation_code
		`).
		Joins("trx_work_order two ON two.work_order_system_number = trx_work_order_allocation.work_order_system_number").
		Where("trx_work_order_allocation.technician_allocation_system_number = ?", techAllocSystemNumber).
		First(&sourceResponseWorkOrderAllocation).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	sourceServiceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(sourceResponseWorkOrderAllocation.ServiceStatusId)
	sourceServiceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(sourceServiceStatusUrl, &sourceServiceStatusPayload, nil); err != nil || sourceServiceStatusPayload.ServiceStatusId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status data"),
		}
	}
	sourceResponseWorkOrderAllocation.ServiceStatusDescription = sourceServiceStatusPayload.ServiceStatusDescription

	for _, description := range serviceStatusList {
		if description == sourceResponseWorkOrderAllocation.ServiceStatusDescription {
			sourceCheckServiceStatus = true
			break
		}
	}
	if !sourceCheckServiceStatus {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("can not moving job allocation, service status the data already " + sourceServiceStatusPayload.ServiceStatusDescription),
		}
	}

	err = tx.Model(&entitiesWorkOrderAllocation).
		Select(`
			trx_work_order_allocation.technician_allocation_system_number,
			trx_work_order_allocation.service_status_id,
			two.vehicle_id,
			trx_work_order_allocation.technician_id,
			trx_work_order_allocation.work_order_system_number,
			trx_work_order_allocation.work_order_document_number,
			trx_work_order_allocation.operation_code
		`).
		Joins("trx_work_order two ON two.work_order_system_number = trx_work_order_allocation.work_order_system_number").
		Where("YEAR(trx_work_order_allocation.tech_alloc_start_date) = ?", req.ServiceDate.Year()).
		Where("MONTH(trx_work_order_allocation.tech_alloc_start_date) = ?", req.ServiceDate.Month()).
		Where("DAY(trx_work_order_allocation.tech_alloc_start_date) = ?", req.ServiceDate.Day()).
		Where("trx_work_order_allocation.technician_id = ?", req.TechnicianId).
		Where("trx_work_order_allocation.sequence_number >= ?", req.SequenceNumber).
		Order("trx_work_order_allocation.sequence_number ASC").
		First(&targetResponseWorkOrderAllocation).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	targetServiceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(targetResponseWorkOrderAllocation.ServiceStatusId)
	targetServiceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(targetServiceStatusUrl, &targetServiceStatusPayload, nil); err != nil || targetServiceStatusPayload.ServiceStatusId == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status data"),
		}
	}
	targetResponseWorkOrderAllocation.ServiceStatusDescription = targetServiceStatusPayload.ServiceStatusDescription

	for _, description := range serviceStatusList {
		if description == targetResponseWorkOrderAllocation.ServiceStatusDescription {
			targetCheckServiceStatus = true
			break
		}
	}
	if !targetCheckServiceStatus {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("can not moving job allocation, service status the data already " + targetResponseWorkOrderAllocation.ServiceStatusDescription),
		}
	}
	// end validate job status

	// begin get work time
	entitiesShiftSchedule := masterentities.ShiftSchedule{}
	shiftTempTable := masterpayloads.ShiftScheduleOutstandingJAResponse{}
	var effectiveDate time.Time
	serviceDate := time.Date(req.ServiceDate.Year(), req.ServiceDate.Month(), req.ServiceDate.Day(), 0, 0, 0, 0, req.ServiceDate.Location())

	err = tx.Model(&entitiesShiftSchedule).
		Select("effective_date").
		Where("company_id = ?", req.CompanyId).
		Where("effective_date <= ?", serviceDate).
		Order("effective_date DESC").
		First(&effectiveDate).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	err = tx.Model(&entitiesShiftSchedule).
		Where("company_id = ?", req.CompanyId).
		Where("effective_date = ?", effectiveDate).
		First(&shiftTempTable).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// end get work

	// begin sequencing target alloc

	if req.TechnicianId == req.OriginalTechnicianId {
		if req.SequenceNumber < req.OriSequenceNumber {
			err = tx.Model(&entitiesWorkOrderAllocation).
				Where("YEAR(tech_alloc_start_date) = ?", serviceDate.Year()).
				Where("MONTH(tech_alloc_start_date) = ?", serviceDate.Month()).
				Where("DAY(tech_alloc_start_date) = ?", serviceDate.Day()).
				Where("technician_id = ?", req.TechnicianId).
				Where("sequence_number >= ?", req.SequenceNumber).
				Update("sequence_number", gorm.Expr("ISNULL(sequence_number, 0) + 1")).Error
		} else {
			err = tx.Model(&entitiesWorkOrderAllocation).
				Where("YEAR(tech_alloc_start_date) = ?", serviceDate.Year()).
				Where("MONTH(tech_alloc_start_date) = ?", serviceDate.Month()).
				Where("DAY(tech_alloc_start_date) = ?", serviceDate.Day()).
				Where("technician_id = ?", req.TechnicianId).
				Where("sequence_number > ?", req.SequenceNumber).
				Update("sequence_number", gorm.Expr("ISNULL(sequence_number, 0) + 1")).Error
		}

		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		err = tx.Model(&entitiesWorkOrderAllocation).
			Where("YEAR(tech_alloc_start_date) = ?", serviceDate.Year()).
			Where("MONTH(tech_alloc_start_date) = ?", serviceDate.Month()).
			Where("DAY(tech_alloc_start_date) = ?", serviceDate.Day()).
			Where("technician_id = ?", req.OriginalTechnicianId).
			Where("sequence_number > ?", req.OriSequenceNumber).
			Update("sequence_number", gorm.Expr("sequence_number - 1")).Error

	} else {
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where("YEAR(tech_alloc_start_date) = ?", serviceDate.Year()).
			Where("MONTH(tech_alloc_start_date) = ?", serviceDate.Month()).
			Where("DAY(tech_alloc_start_date) = ?", serviceDate.Day()).
			Where("technician_id = ?", req.TechnicianId).
			Where("sequence_number >= ?", req.SequenceNumber).
			Update("sequence_number", gorm.Expr("ISNULL(sequence_number, 0) + 1")).Error
	}

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// end sequencing target alloc

	// BEGIN moving alloc to target

	companyRefUrl := config.EnvConfigs.GeneralServiceUrl + "company-reference/" + strconv.Itoa(req.CompanyId)
	companyRefPayload := transactionjpcbpayloads.OutstandingJACompanyRefPayload{}
	if err := utils.Get(companyRefUrl, &companyRefPayload, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching company reference data"),
		}
	}

	var currentDate time.Time
	var currentTime float64
	var targetTechAllocLastEndTime float64

	currentDate = time.Now().Add(time.Duration(companyRefPayload.TimeDifference) * time.Hour)

	if currentDate.Year() == serviceDate.Year() &&
		currentDate.Month() == serviceDate.Month() &&
		currentDate.Day() == serviceDate.Day() &&
		req.SequenceNumber == 1 {
		currentTime = utils.TimeValue(currentDate)
		targetTechAllocLastEndTime = currentTime
	} else {
		err = tx.Model(&entitiesWorkOrderAllocation).
			Select("tech_alloc_last_end_time").
			Where("YEAR(tech_alloc_start_date) = ?", serviceDate.Year()).
			Where("MONTH(tech_alloc_start_date) = ?", serviceDate.Month()).
			Where("DAY(tech_alloc_start_date) = ?", serviceDate.Day()).
			Where("technician_id = ?", req.TechnicianId).
			Where("sequence_number < ?", req.SequenceNumber).
			Order("sequence_number DESC").
			First(&targetTechAllocLastEndTime).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// begin factor_x

	var qcExtraFrt float64
	var reorderJob bool
	var isExpress bool
	var frtQuantity float64
	var frtQuantityOri float64
	factorX := 1.0

	employeeUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(req.TechnicianId)
	employeePayload := transactionjpcbpayloads.OutstandingJAEmployeePayload{}
	if err := utils.Get(employeeUrl, &employeePayload, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching employee / technician data"),
		}
	}
	if employeePayload.FactorX > 1 {
		factorX = employeePayload.FactorX
	}

	var checkExistWorkOrder int
	err = tx.Model(&entitiesWorkOrderAllocation).
		Select("ISNULL(work_order_system_number, 0) AS work_order_system_number").
		Where("technician_allocation_system_number = ?", techAllocSystemNumber).
		Pluck("work_order_system_number", &checkExistWorkOrder).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if checkExistWorkOrder != 0 {
		fetchFrtPayload := transactionjpcbpayloads.OutstandingJobAllocationFetchFRTPayload{}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Select(`
				twod.quality_control_extra_frt,
				trx_work_order_allocation.re_order,
				trx_work_order_allocation.is_express,
				trx_work_order_allocation.tech_alloc_total_time,
				trx_work_order_allocation.frt
			`).
			Joins("INNER JOIN trx_work_order two ON two.work_order_system_number = trx_work_order_allocation.work_order_system_number").
			Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = two.work_order_system_number").
			Where("trx_work_order_allocation.technician_allocation_system_number = ?", techAllocSystemNumber).
			First(&fetchFrtPayload).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		qcExtraFrt = fetchFrtPayload.QualityControlExtraFrt
		reorderJob = fetchFrtPayload.ReOrder
		isExpress = fetchFrtPayload.IsExpress
		frtQuantity = fetchFrtPayload.TechAllocTotalTime
		frtQuantityOri = fetchFrtPayload.Frt

	} else {
		fetchFrtPayload := transactionjpcbpayloads.OutstandingJobAllocationFetchFRTPayload{}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: techAllocSystemNumber}).
			First(&fetchFrtPayload).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		isExpress = fetchFrtPayload.IsExpress
		frtQuantity = fetchFrtPayload.TechAllocTotalTime
		frtQuantityOri = fetchFrtPayload.Frt
	}

	if reorderJob {
		frtQuantity = qcExtraFrt
	} else if !isExpress {
		frtQuantity = frtQuantityOri * factorX
	}

	err = tx.Model(&entitiesWorkOrderAllocation).
		Where(transactionworkshopentities.WorkOrderAllocation{
			TechAllocSystemNumber: techAllocSystemNumber,
		}).
		Updates(transactionworkshopentities.WorkOrderAllocation{
			TechAllocTotalTime: frtQuantity,
			FactorX:            factorX,
		}).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// end factor_x

	workOrderAllocationPayload := transactionworkshopentities.WorkOrderAllocation{}
	err = tx.Model(&entitiesWorkOrderAllocation).
		Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: techAllocSystemNumber}).
		First(&workOrderAllocationPayload).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var techAllocLastStartTime float64
	var techAllocLastEndTime float64

	var startTime float64
	if targetTechAllocLastEndTime == 0 {
		startTime = shiftTempTable.StartTime
	} else {
		startTime = targetTechAllocLastEndTime
	}
	techAllocLastStartTime = startTime

	startTimeAdded := startTime + workOrderAllocationPayload.TechAllocTotalTime

	if startTime <= shiftTempTable.RestStartTime && startTimeAdded >= shiftTempTable.RestStartTime {
		techAllocLastEndTime = shiftTempTable.RestEndTime + startTime + workOrderAllocationPayload.TechAllocTotalTime - shiftTempTable.RestStartTime
	} else {
		techAllocLastEndTime = startTime + workOrderAllocationPayload.TechAllocTotalTime
	}

	serviceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status-by-description/Draft"
	serviceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayload, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status draft data"),
		}
	}

	err = tx.Model(&entitiesWorkOrderAllocation).
		Where(transactionworkshopentities.WorkOrderAllocation{
			TechAllocSystemNumber: techAllocSystemNumber,
		}).
		Updates(transactionworkshopentities.WorkOrderAllocation{
			TechnicianId:           req.TechnicianId,
			SequenceNumber:         req.SequenceNumber,
			TechAllocLastStartTime: techAllocLastStartTime,
			TechAllocLastEndTime:   techAllocLastEndTime,
			ServStatus:             serviceStatusPayload.ServiceStatusId,
		}).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var checkServiceLogExist int
	entitiesServiceLog := transactionworkshopentities.ServiceLog{}
	err = tx.Model(&entitiesServiceLog).
		Select("COUNT(*) AS count").
		Where(transactionworkshopentities.ServiceLog{TechnicianAllocationSystemNumber: techAllocSystemNumber}).
		Pluck("count", &checkServiceLogExist).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if checkServiceLogExist > 0 {
		err = tx.Delete(&entitiesWorkOrderAllocation, techAllocSystemNumber).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		insertServiceLogPayload := transactionjpcbpayloads.OutstandingJobAllocationInsertServiceLogPayload{}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Select(`
				trx_work_order_allocation.company_id,
				trx_work_order_allocation.technician_allocation_system_number,
				1 AS technician_allocation_line,
				two.work_order_system_number,
				two.work_order_document_number,
				two.work_order_date,
				trx_work_order_allocation.operation_code,
				trx_work_order_allocation.technician_id,
				trx_work_order_allocation.shift_code,
				twod.frt_quantity,
				trx_work_order_allocation.tech_alloc_last_start_date,
				trx_work_order_allocation.tech_alloc_last_end_date,
				trx_work_order_allocation.tech_alloc_last_start_time,
				trx_work_order_allocation.sequence_number
			`).
			Joins("INNER JOIN trx_work_order two ON two.work_order_system_number = trx_work_order_allocation.work_order_system_number").
			Joins("INNER JOIN trx_work_order_detail twod ON twod.work_order_system_number = two.work_order_system_number").
			Where("trx_work_order_allocation.technician_allocation_system_number = ?", techAllocSystemNumber).
			First(&insertServiceLogPayload).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		entitiesServiceLog = transactionworkshopentities.ServiceLog{
			CompanyId:                        insertServiceLogPayload.CompanyId,
			TechnicianAllocationSystemNumber: insertServiceLogPayload.TechnicianAllocationSystemNumber,
			TechnicianAllocationLine:         insertServiceLogPayload.TechnicianAllocationLine,
			WorkOrderSystemNumber:            insertServiceLogPayload.WorkOrderSystemNumber,
			WorkOrderDocumentNumber:          insertServiceLogPayload.WorkOrderDocumentNumber,
			WorkOrderDate:                    insertServiceLogPayload.WorkOrderDate.Format("2006-01-02 15:04:05"),
			OperationItemCode:                insertServiceLogPayload.OperationCode,
			TechnicianId:                     insertServiceLogPayload.TechnicianId,
			ShiftCode:                        insertServiceLogPayload.ShiftCode,
			Frt:                              insertServiceLogPayload.FrtQuantity,
			ServiceStatusId:                  serviceStatusPayload.ServiceStatusId,
			StartDatetime:                    insertServiceLogPayload.TechAllocLastStartDate,
			EndDatetime:                      insertServiceLogPayload.TechAllocLastEndDate,
			ActualTime:                       0,
			EstimatedPendingTime:             0,
			PendingTime:                      0,
			ServiceReasonId:                  0,
			EmpGroupId:                       0,
			Remark:                           "",
			ActualStartTime:                  insertServiceLogPayload.TechAllocLastStartTime,
			SequenceNumber:                   insertServiceLogPayload.SequenceNumber,
			ServiceLogSystemNumber:           0,
			WorkOrderOperationId:             0,
			ShiftScheduleId:                  0,
			WorkOrderLine:                    0,
		}

		err = tx.Save(&entitiesServiceLog).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// END moving alloc to target alloc

	// begin resequence source alloc

	if req.TechnicianId != req.OriginalTechnicianId {
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where("YEAR(tech_alloc_start_date) = ?", req.ServiceDate.Year()).
			Where("MONTH(tech_alloc_start_date) = ?", req.ServiceDate.Month()).
			Where("DAY(tech_alloc_start_date) = ?", req.ServiceDate.Day()).
			Where("technician_id = ?", req.OriginalTechnicianId).
			Where("sequence_number > ?", req.SequenceNumber).
			Update("sequence_number", gorm.Expr("sequence_number - 1")).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	var sourceFirstTechAllocSystenNumber int
	var sourceFirstServiceStatus int
	var sourceFirstSequenceNumber int

	sourceFirstTechAllocation := transactionworkshopentities.WorkOrderAllocation{}
	err = tx.Model(&entitiesWorkOrderAllocation).
		Where("YEAR(tech_alloc_start_date) = ?", req.ServiceDate.Year()).
		Where("MONTH(tech_alloc_start_date) = ?", req.ServiceDate.Month()).
		Where("DAY(tech_alloc_start_date) = ?", req.ServiceDate.Day()).
		Where("technician_id = ?", req.OriginalTechnicianId).
		Order("sequence_number ASC").
		First(&sourceFirstTechAllocation).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	sourceFirstTechAllocSystenNumber = sourceFirstTechAllocation.TechAllocSystemNumber
	sourceFirstServiceStatus = sourceFirstTechAllocation.ServStatus
	sourceFirstSequenceNumber = sourceFirstTechAllocation.SequenceNumber

	companyRefUrl = config.EnvConfigs.GeneralServiceUrl + "company-reference/" + strconv.Itoa(req.CompanyId)
	companyRefPayload = transactionjpcbpayloads.OutstandingJACompanyRefPayload{}
	if err := utils.Get(companyRefUrl, &companyRefPayload, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching company reference data"),
		}
	}

	serviceStatusUrl = config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(sourceFirstServiceStatus)
	serviceStatusPayload = transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayload, nil); err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error fetching service status data"),
		}
	}

	checkServiceStatus := false
	for _, description := range serviceStatusList {
		if description == serviceStatusPayload.ServiceStatusDescription {
			checkServiceStatus = true
			break
		}
	}

	// move source first alloc Seq
	currentDate = time.Now().Add(time.Duration(companyRefPayload.TimeDifference) * time.Hour)
	currentTime = utils.TimeValue(currentDate)

	techAllocLastStartTime = 0.0
	techAllocLastEndTime = 0.0
	startTime = 0.0

	if currentDate.Year() == serviceDate.Year() &&
		currentDate.Month() == serviceDate.Month() &&
		currentDate.Day() == serviceDate.Day() &&
		sourceFirstSequenceNumber == 1 &&
		checkServiceStatus {

		if currentTime < shiftTempTable.StartTime {
			startTime = shiftTempTable.StartTime
		} else {
			startTime = currentTime
		}

		techAllocLastStartTime = startTime
		if startTime < shiftTempTable.RestStartTime && (startTime+sourceFirstTechAllocation.TechAllocTotalTime) >= shiftTempTable.RestStartTime {
			techAllocLastEndTime = shiftTempTable.RestEndTime + startTime + sourceFirstTechAllocation.TechAllocTotalTime - shiftTempTable.RestStartTime
		} else {
			techAllocLastEndTime = startTime + sourceFirstTechAllocation.TechAllocTotalTime
		}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where(transactionworkshopentities.WorkOrderAllocation{
				TechAllocSystemNumber: sourceFirstTechAllocSystenNumber,
			}).
			Updates(transactionworkshopentities.WorkOrderAllocation{
				TechAllocLastStartTime: techAllocLastStartTime,
				TechAllocLastEndTime:   techAllocLastEndTime,
			}).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	} else if checkServiceStatus {
		techAllocLastStartTime = shiftTempTable.StartTime
		if shiftTempTable.StartTime+sourceFirstTechAllocation.TechAllocTotalTime == shiftTempTable.RestStartTime {
			techAllocLastEndTime = shiftTempTable.RestEndTime + shiftTempTable.RestStartTime + sourceFirstTechAllocation.TechAllocTotalTime - shiftTempTable.RestStartTime
		} else {
			techAllocLastEndTime = shiftTempTable.StartTime + sourceFirstTechAllocation.TechAllocTotalTime
		}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where(transactionworkshopentities.WorkOrderAllocation{
				TechAllocSystemNumber: sourceFirstTechAllocSystenNumber,
			}).
			Updates(transactionworkshopentities.WorkOrderAllocation{
				TechAllocLastStartTime: techAllocLastStartTime,
				TechAllocLastEndTime:   techAllocLastEndTime,
			}).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		if workOrderAllocationPayload.BookingSystemNumber > 0 {
			var serviceTime float32
			entitiesBookingEstimation := transactionworkshopentities.BookingEstimation{}
			err = tx.Model(&entitiesBookingEstimation).
				Select("booking_service_time").
				Joins("INNER JOIN trx_booking_estimation_allocation tbea ON tbea.booking_system_number = trx_booking_estimation.booking_system_number").
				Joins("INNER JOIN trx_work_order_allocation twoa ON twoa.booking_system_number = trx_booking_estimation.booking_system_number").
				Where("trx_work_order_allocation.technician_allocation_system_number = ?", techAllocSystemNumber).
				Pluck("booking_service_time", &serviceTime).Error
			if err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			bookingServiceTime := float64(serviceTime)

			techAllocLastStartTime = bookingServiceTime
			if bookingServiceTime+workOrderAllocationPayload.TechAllocTotalTime <= bookingServiceTime && bookingServiceTime+workOrderAllocationPayload.TechAllocTotalTime >= shiftTempTable.RestStartTime {
				techAllocLastEndTime = shiftTempTable.RestEndTime + bookingServiceTime + workOrderAllocationPayload.TechAllocTotalTime - shiftTempTable.RestStartTime
			} else {
				techAllocLastEndTime = bookingServiceTime + workOrderAllocationPayload.TechAllocTotalTime
			}
			err = tx.Model(&entitiesWorkOrderAllocation).
				Where(transactionworkshopentities.WorkOrderAllocation{
					TechAllocSystemNumber: techAllocSystemNumber,
				}).
				Updates(transactionworkshopentities.WorkOrderAllocation{
					TechAllocLastStartTime: techAllocLastStartTime,
					TechAllocLastEndTime:   techAllocLastEndTime,
				}).Error
			if err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}

			techAllocLastStartTime = bookingServiceTime
			if bookingServiceTime+sourceFirstTechAllocation.TechAllocTotalTime <= bookingServiceTime && bookingServiceTime+sourceFirstTechAllocation.TechAllocTotalTime >= shiftTempTable.RestStartTime {
				techAllocLastEndTime = shiftTempTable.RestEndTime + bookingServiceTime + sourceFirstTechAllocation.TechAllocTotalTime - shiftTempTable.RestStartTime
			} else {
				techAllocLastEndTime = bookingServiceTime + sourceFirstTechAllocation.TechAllocTotalTime
			}
			err = tx.Model(&entitiesWorkOrderAllocation).
				Where(transactionworkshopentities.WorkOrderAllocation{
					TechAllocSystemNumber: sourceFirstTechAllocSystenNumber,
				}).
				Updates(transactionworkshopentities.WorkOrderAllocation{
					TechAllocLastStartTime: techAllocLastStartTime,
					TechAllocLastEndTime:   techAllocLastEndTime,
				}).Error
			if err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
		}
	}

	updateServiceLogTemplate := []transactionworkshopentities.WorkOrderAllocation{}
	err = tx.Model(&entitiesWorkOrderAllocation).
		Where("YEAR(tech_alloc_start_date) = ?", req.ServiceDate.Year()).
		Where("MONTH(tech_alloc_start_date) = ?", req.ServiceDate.Month()).
		Where("DAY(tech_alloc_start_date) = ?", req.ServiceDate.Day()).
		Where("technician_id IN (?, ?)", req.TechnicianId, req.OriginalTechnicianId).
		Scan(&updateServiceLogTemplate).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// end resequencing alloc

	for _, template := range updateServiceLogTemplate {
		err = tx.Model(&entitiesServiceLog).
			Where(transactionworkshopentities.ServiceLog{
				TechnicianAllocationSystemNumber: template.TechAllocSystemNumber,
				TechnicianId:                     template.TechnicianId,
			}).
			Update("sequence_number", template.SequenceNumber).Error
		if err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	response = transactionjpcbpayloads.OutstandingJobAllocationUpdateResponse{
		SourceFirstTechAllocSystenNumber: sourceFirstTechAllocSystenNumber,
		TechAllocSystemNumber:            techAllocSystemNumber,
	}

	return response, nil
}

// uspg_atWoTechAlloc_ReCalCulateTimeJob
// IF @Option = 0
func (r *OutstandingJobAllocationRepositoryImpl) ReCalculateTimeJob(tx *gorm.DB, techAllocSystemNumber int) *exceptions.BaseErrorResponse {
	serviceStatusUrl := config.EnvConfigs.GeneralServiceUrl + "service-status-by-description/Draft"
	serviceStatusPayload := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayload, nil); err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entitiesWorkOrderAllocation := transactionworkshopentities.WorkOrderAllocation{}
	responseWorkOrderAllocation := transactionworkshopentities.WorkOrderAllocation{}
	err := tx.Model(&entitiesWorkOrderAllocation).
		Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: techAllocSystemNumber}).
		First(&responseWorkOrderAllocation).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	bookingSystemNumber := responseWorkOrderAllocation.BookingSystemNumber
	techAllocStartTime := responseWorkOrderAllocation.TechAllocLastStartTime
	bookingServiceTime := responseWorkOrderAllocation.BookingServiceTime
	techAllocLastStartTime := responseWorkOrderAllocation.TechAllocLastStartTime
	frtQuantity := responseWorkOrderAllocation.TechAllocTotalTime
	techAllocLastEndTime := responseWorkOrderAllocation.TechAllocLastEndTime
	techAllocStartDate := responseWorkOrderAllocation.TechAllocStartDate
	technicianId := responseWorkOrderAllocation.TechnicianId
	sequenceNumber := responseWorkOrderAllocation.SequenceNumber
	companyId := responseWorkOrderAllocation.CompanyId
	shiftCode := responseWorkOrderAllocation.ShiftCode
	factorX := responseWorkOrderAllocation.FactorX
	if factorX < 1 {
		factorX = 1
	}

	if bookingSystemNumber > 0 && techAllocStartTime < bookingServiceTime {
		techAllocStartTime = bookingServiceTime
	}

	entitiesShiftSchedule := masterentities.ShiftSchedule{}
	responseShiftSchedule := masterentities.ShiftSchedule{}
	err = tx.Model(&entitiesShiftSchedule).
		Where(masterentities.ShiftSchedule{
			CompanyId: companyId,
			ShiftCode: shiftCode,
		}).
		Order("effective_date DESC").
		First(&responseShiftSchedule).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	restStartTime := responseShiftSchedule.RestStartTime
	restEndTime := responseShiftSchedule.RestEndTime

	if techAllocStartTime > restStartTime && techAllocLastStartTime < restEndTime {
		techAllocStartTime = restEndTime
		techAllocLastEndTime = restEndTime + frtQuantity
	} else if techAllocLastStartTime < restStartTime && techAllocLastEndTime > restStartTime && techAllocLastEndTime < restEndTime {
		techAllocLastEndTime = restEndTime + (frtQuantity - (restStartTime - techAllocLastStartTime))
	}

	if serviceStatusPayload.ServiceStatusDescription == "DRAFT" {
		entitiesWorkOrderAllocation := transactionworkshopentities.WorkOrderAllocation{}
		err = tx.Model(&entitiesWorkOrderAllocation).
			Where(transactionworkshopentities.WorkOrderAllocation{TechAllocSystemNumber: techAllocSystemNumber}).
			Updates(transactionworkshopentities.WorkOrderAllocation{
				TechAllocLastStartTime: techAllocLastStartTime,
				TechAllocLastEndTime:   techAllocLastEndTime,
			}).Error
		if err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	startTime := techAllocLastEndTime

	newSequenceWOAlloc := transactionworkshopentities.WorkOrderAllocation{}
	err = tx.Model(&entitiesWorkOrderAllocation).
		Where("company_id = ?", companyId).
		Where("tech_alloc_start_date = ?", techAllocStartDate).
		Where("technician_id = ?", technicianId).
		Where("sequence_number > ?", sequenceNumber).
		Order("sequence_number ASC").
		First(&newSequenceWOAlloc).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	serviceStatusUrl = config.EnvConfigs.GeneralServiceUrl + "service-status/" + strconv.Itoa(newSequenceWOAlloc.ServStatus)
	serviceStatusPayloads := transactionjpcbpayloads.ServiceStatusPayload{}
	if err := utils.Get(serviceStatusUrl, &serviceStatusPayloads, nil); err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var endTime float64

	if newSequenceWOAlloc.TechAllocSystemNumber > 0 {
		newSeqFrtQuantity := newSequenceWOAlloc.TechAllocTotalTime * factorX

		if newSequenceWOAlloc.BookingSystemNumber > 0 && startTime < newSequenceWOAlloc.BookingServiceTime {
			startTime = newSequenceWOAlloc.BookingServiceTime
		}

		if startTime > restStartTime {
			endTime = startTime + newSeqFrtQuantity
		} else if startTime < restStartTime && (startTime+newSeqFrtQuantity) > restStartTime {
			endTime = restEndTime + ((startTime + newSeqFrtQuantity) - restStartTime)
		} else {
			endTime = startTime + newSeqFrtQuantity
		}

		if endTime > 21.00 {
			endTime = 21.00
		}

		if serviceStatusPayloads.ServiceStatusDescription == "Draft" {
			err = tx.Model(&entitiesWorkOrderAllocation).
				Where(transactionworkshopentities.WorkOrderAllocation{
					TechAllocSystemNumber: techAllocSystemNumber,
				}).
				Updates(transactionworkshopentities.WorkOrderAllocation{
					TechAllocStartTime:     startTime,
					TechAllocEndTime:       endTime,
					TechAllocLastStartTime: startTime,
					TechAllocLastEndTime:   endTime,
				}).Error
			if err != nil {
				return &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
		}
	}

	return nil
}
